package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sentinel-dashboard/db"
	"github.com/sentinel-dashboard/handler"
)

type Sentinels struct {
	ID   int
	Name string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbConn := db.New(os.Getenv("DB_FILE_NAME"))
	defer dbConn.Close()

	fmt.Println("Run DB Migration")
	dbConn.Migrate()

	fmt.Println("Starting Sentinel Manager Server")
	h := handler.New(dbConn)
	r := h.Router()
	r.Run(":" + os.Getenv("BACKEND_PORT"))
}
