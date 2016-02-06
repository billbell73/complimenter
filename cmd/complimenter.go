package main

import (
	"html/template"
	"database/sql"
	"net/http"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type compliment struct {
	Body string
}

func fetchMaxId(db *sql.DB) int {
  var maxId int
  row := db.QueryRow("SELECT MAX(id) FROM compliments")
  err := row.Scan(&maxId)
  checkErr(err)

	return maxId
}

func fetchCompliment(db *sql.DB, id int) compliment {
	var body string
  row := db.QueryRow("SELECT body FROM compliments WHERE id=?", id)
  err := row.Scan(&body)
  checkErr(err)

  return compliment{body}
}

func randomId(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n) + 1
}

func showHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "/test")
	checkErr(err)
	defer db.Close()

	err = db.Ping()
	checkErr(err)
	maxId := fetchMaxId(db)
	chosen_id := randomId(maxId)
	compliment := fetchCompliment(db, chosen_id)

	t, _ := template.ParseFiles("../views/show.html")
  t.Execute(w, compliment)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("../views/index.html")
  t.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/show", showHandler)
	http.ListenAndServe(":8080", nil)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
