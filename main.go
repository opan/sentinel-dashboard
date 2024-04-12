package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sentinel-dashboard/db"
	"github.com/sentinel-dashboard/handler"
)

func main() {
	ctx := context.TODO()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s\n", err)
	}

	dbConn := db.New(os.Getenv("DB_FILE_NAME"))
	defer dbConn.Close()

	fmt.Println("Run DB Migration")
	dbConn.Migrate()

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fmt.Println("Starting Sentinel Manager Server")
	h := handler.New(dbConn)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("BACKEND_PORT")),
		Handler: h.Router(),
	}

	// init the server in goroutine so that it won't block the graceful shutdown handling below
	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// listen for the interrupt signal
	<-ctx.Done()

	// restore default behaviour on the interrupt signal and notify user of shutdown
	stop()
	log.Println("shutting down gracefully, press ctrl+C to force shutdown")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting. Byebye!")
}
