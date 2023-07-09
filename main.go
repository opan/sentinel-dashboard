package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Sentinels struct {
	ID   int
	Name string
}

func main() {
	// handler.Router()

	db, err := sql.Open("sqlite3", "./sentinel.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	sqlQuery := `
	CREATE TABLE sentinels (
		id INTEGER NOT NULL PRIMARY KEY,
		name TEXT NOT NULL
	)`

	_, err = db.Exec(sqlQuery)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlQuery)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("INSERT INTO sentinels (name) VALUES (?)")
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	_, err = stmt.Exec("main-cluster")
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT * FROM sentinels")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var name string

		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(id, name)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
