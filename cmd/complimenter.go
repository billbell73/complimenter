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

func fetchCompliments() []string {
	db, err := sql.Open("mysql", "/test")
	checkErr(err)
	defer db.Close()

	err = db.Ping()
	checkErr(err)

	rows, err := db.Query("SELECT * FROM compliments")
	checkErr(err)

	var compliments []string

	for rows.Next() {
		var body string
		var id int
		err = rows.Scan(&id, &body)
		checkErr(err)
		compliments = append(compliments, body)
	}
	return compliments
}

func randomIntLessThan(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

func main() {
	compliments := fetchCompliments()
	index := randomIntLessThan(len(compliments))
	fmt.Println(compliments[index])
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
