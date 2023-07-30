package main

import (
	"fmt"

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

	fmt.Println("Run DB Migration")
	dbConn.Migrate()

	fmt.Println("Connecting to Redis Sentinel Servers")
	s := redis.NewSentinelClient(&redis.Options{
		Addr: ":26379",
	})

	h := handler.New(dbConn, s)
	h.Router()
}
