package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
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

	dbConn.Migrate()

	s := redis.NewSentinelClient(&redis.Options{
		Addr: ":26379",
	})

	h := handler.New(dbConn, s)
	h.Router()
}
