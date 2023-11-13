package handler_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sentinel-dashboard/handler"
	"github.com/stretchr/testify/assert"
)

func Test_registerSentinelHandler(t *testing.T) {
	dbConn := setupTest()
	defer dbConn.Close()

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
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, `{"errors":[],"msg":"Sentinel successfully registered"}`, string(respData))
}
