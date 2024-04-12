package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/joho/godotenv"
	"github.com/sentinel-manager/db"
	"github.com/sentinel-manager/handler"
)

func setupTest() db.DB {
	_ = godotenv.Load(".env.test")
	dbConn := db.New(os.Getenv("DB_FILE_NAME"))
	dbConn.Migrate()
	return dbConn
}

func makeRequest(dbConn db.DB, method, url string, body interface{}) *httptest.ResponseRecorder {
	reqBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.Header.Add("Content-Type", "application/json")

	h := handler.New(dbConn)
	r := h.Router()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}
