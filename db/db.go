package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type database struct {
	dbConn *sql.DB
}

type DB interface {
	Close()
	Migrate()
	Drop()
}

func CreateConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "./sentinels.db")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func (d *database) Close() {
	d.dbConn.Close()
}

func (d *database) Migrate() {
	sqlQuery := `
		CREATE TABLE IF NOT EXISTS sentinels (
		id INTEGER NOT NULL PRIMARY KEY,
		name TEXT NOT NULL,
		hosts TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP;
	`

	_, err := d.dbConn.Exec(sqlQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *database) Drop() {
	os.Remove("./sentinel.db")
}

func New() DB {
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
