package main

import (
	"database/sql"
	"html/template"
	"math/rand"
	"net/http"
	"time"
	"os"
	"log"

	_ "github.com/billbell73/complimenter/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
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

func showHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	maxId := fetchMaxId(db)
	chosen_id := randomId(maxId)
	compliment := fetchCompliment(db, chosen_id)

	t, _ := template.ParseFiles("../views/show.html")
	t.Execute(w, compliment)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/index.html")
	checkErr(err)
	t.Execute(w, nil)
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("../views/new.html")
	t.Execute(w, nil)
}

func saveHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	body := r.FormValue("body")

	stmt, err := db.Prepare("INSERT INTO compliments(body) VALUES(?)")
	checkErr(err)

	_, err = stmt.Exec(body)
	checkErr(err)

	http.Redirect(w, r, "/", http.StatusFound)
}

func makeDbHandler(fn func(http.ResponseWriter, *http.Request, *sql.DB)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
		checkErr(err)
		defer db.Close()

		err = db.Ping()
		checkErr(err)
		fn(w, r, db)
	}
}

func serveSingle(pattern string, filename string) {
    http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, filename)
    })
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/new", newHandler)
	http.HandleFunc("/show", makeDbHandler(showHandler))
	http.HandleFunc("/save", makeDbHandler(saveHandler))

	serveSingle("/sitemap.xml", "./sitemap.xml")
    serveSingle("/favicon.ico", "./favicon.ico")
    serveSingle("/robots.txt", "./robots.txt")


	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.ListenAndServe(":"+port, nil)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
