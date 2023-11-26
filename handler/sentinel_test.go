package handler_test

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"github.com/sentinel-dashboard/db"
	"github.com/sentinel-dashboard/model"
	"github.com/stretchr/testify/assert"
)

func Test_registerSentinelHandler(t *testing.T) {
	dbConn := setupTest()
	defer dbConn.Close()

	newSentinel := model.Sentinel{
		Name:  "sentinel-test",
		Hosts: "10.23.22.10:26379,10.23.22.11:26379",
	}

	w := makeRequest(dbConn, "POST", "/sentinel/register", newSentinel)

	var res map[string]string
	json.Unmarshal(w.Body.Bytes(), &res)
	_, exists := res["msg"]

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, true, exists)
}

func Test_getSentinelHandler(t *testing.T) {
	dbConn := setupTest()
	defer dbConn.Close()

	setupDummyDataSentinelHandler(dbConn)

	w := makeRequest(dbConn, "GET", "/sentinel/1", nil)

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_getAllSentinelHandler(t *testing.T) {
	dbConn := setupTest()
	defer dbConn.Close()

	setupDummyDataSentinelHandler(dbConn)

	w := makeRequest(dbConn, "GET", "/sentinel", nil)

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, w.Code)
	// var sentinels []model.Sentinel

	// json.Unmarshal(byte(res["data"]), &sentinels)
	// fmt.Println(len(sentinels))
}

func setupDummyDataSentinelHandler(dbConn db.DB) {
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
