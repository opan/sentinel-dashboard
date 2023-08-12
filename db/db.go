package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type database struct {
	db *sql.DB
}

type DB interface {
	GetConnection() *sql.DB
	Close()
	Migrate()
	Drop()
}

func (d *database) GetConnection() *sql.DB {
	return d.db
}

func (d *database) Close() {
	d.db.Close()
}

func (d *database) Migrate() {
	sqlQuery := `
		CREATE TABLE IF NOT EXISTS sentinels (
		id INTEGER NOT NULL PRIMARY KEY,
		name TEXT NOT NULL,
		hosts TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP);
	`

	_, err := d.db.Exec(sqlQuery)
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
		db: db,
	}
}
