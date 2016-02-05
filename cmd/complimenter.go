package main

import (
	"database/sql"
	"fmt"
	// "net/http"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type compliment struct {
	body string
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

func main() {
	db, err := sql.Open("mysql", "/test")
	checkErr(err)
	defer db.Close()

	err = db.Ping()
	checkErr(err)

	maxId := fetchMaxId(db)
	chosen_id := randomId(maxId)
	compliment := fetchCompliment(db, chosen_id)

	fmt.Println(compliment)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
