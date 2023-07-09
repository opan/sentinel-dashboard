package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type database struct {
	dbConn *sql.DB
}

type DB interface {
	CloseConn()
}

func CreateConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "./sentinels.db")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func (d *database) CloseConn() {
	d.dbConn.Close()
}

func NewConn() DB {
	driver := "sqlite3"
	dbFile := "./sentinels.db"
	db, err := sql.Open(driver, dbFile)
	if err != nil {
		log.Fatal(err)
	}

	return &database{
		dbConn: db,
	}
}
