package handler_test

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sentinel-dashboard/db"
)

func setupTest() db.DB {
	_ = godotenv.Load(".env.test")
	dbConn := db.New(os.Getenv("DB_FILE_NAME"))
	dbConn.Migrate()
	return dbConn
}
