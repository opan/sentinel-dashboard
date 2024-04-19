package db

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

type database struct {
	db *sqlx.DB
}

type DB interface {
	GetConnection() *sqlx.DB
	Close()
	Migrate()
	Drop()
}

func (d *database) GetConnection() *sqlx.DB {
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

		CREATE TABLE IF NOT EXISTS sentinel_masters (
			id INTEGER NOT NULL PRIMARY KEY,
			sentinel_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			ip TEXT NOT NULL,
			port TEXT NOT NULL,
			quorum TEXT NOT NULL,
			options TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (sentinel_id)
				REFERENCES sentinels
		)
	`

	_, err := d.db.Exec(sqlQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *database) Drop() {
	os.Remove(os.Getenv("DB_FILE_NAME"))
}

func New(dbName string) DB {
	driver := "sqlite3"
	db, err := sqlx.Open(driver, dbName)
	if err != nil {
		log.Fatal(err)
	}

	return &database{
		db: db,
	}
}
