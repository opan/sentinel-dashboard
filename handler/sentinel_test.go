package handler_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sentinel-dashboard/db"
	"github.com/sentinel-dashboard/handler"
)

func Test_registerSentinelHandler(t *testing.T) {
	t.Setenv("DB_FILE_NAME", "./sentinel_test.db")
	dbConn := db.New()
	defer dbConn.Close()
	dbConn.Migrate()

	reqBody := []byte(`{
		"name": "sentinel-test",
		"hosts": "10.23.22.10:26379,10.23.22.11:26379"
	}`)

	req, _ := http.NewRequest("POST", "/sentinel/register", bytes.NewBuffer(reqBody))
	req.Header.Add("Content-Type", "application/json")

	h := handler.New(dbConn)
	r := h.Router()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	respData, _ := io.ReadAll(w.Body)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"errors":[],"msg":"Sentinel successfully register"}`, respData)
}
