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

func randomRead() string {
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

	rand.Seed(time.Now().UnixNano())
	return compliments[rand.Intn(len(compliments))]
}

func main() {
	fmt.Println(randomRead())
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
