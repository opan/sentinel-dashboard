package handler_test

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sentinel-dashboard/db"
	"github.com/sentinel-dashboard/handler"
	"github.com/sentinel-dashboard/model"
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

func Test_getSentinelHandler(t *testing.T) {
	dbConn := setupTest()
	defer dbConn.Close()

	setupDummySentinelHandler(dbConn)

	req, _ := http.NewRequest("GET", "/sentinel", nil)
	h := handler.New(dbConn)
	r := h.Router()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	respData, _ := io.ReadAll(w.Body)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"errors":[],"msg":"","data":[]}`, string(respData))

	// db := dbConn.GetConnection()
	// rows, err := db.Query("SELECT COUNT(id) FROM sentinels")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer rows.Close()
	// var i int
	// for rows.Next() {

	// 	_ = rows.Scan(&i)
	// }

	// fmt.Println(i)

}

func setupDummySentinelHandler(dbConn db.DB) {
	db := dbConn.GetConnection()
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(`INSERT INTO sentinels (name, hosts) VALUES (?, ?)`)
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	sentinelTables := []model.Sentinel{
		{Name: "sentinel-dummy-1", Hosts: "10.12.1.1:26379"},
		{Name: "sentinel-dummy-2", Hosts: "10.12.1.1:26379"},
	}

	for _, v := range sentinelTables {
		_, err = stmt.Exec(v.Name, v.Hosts)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

}
