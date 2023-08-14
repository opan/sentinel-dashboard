package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sentinel-dashboard/db"
	"github.com/sentinel-dashboard/handler"
)

type Sentinels struct {
	ID   int
	Name string
}

func main() {
	dbConn := db.New()
	defer dbConn.Close()

	fmt.Println("Run DB Migration")
	dbConn.Migrate()

	fmt.Println("Starting Sentinel Manager Server")
	h := handler.New(dbConn)
	r := h.Router()
	r.Run(":8282")
}
