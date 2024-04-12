package handler_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/sentinel-manager/db"
	"github.com/sentinel-manager/model"
	"github.com/stretchr/testify/assert"
)

var res struct {
	Msg    string
	Data   []model.Sentinel
	Errors []string
}

func Test_createSentinelHandler(t *testing.T) {
	dbConn := setupTest()
	defer dbConn.Close()

	newSentinel := model.Sentinel{
		Name:  "sentinel-test",
		Hosts: "10.23.22.10:26379,10.23.22.11:26379",
	}

	w := makeRequest(dbConn, "POST", "/sentinel", newSentinel)

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	_, exists := res["msg"]

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, true, exists)
	assert.Equal(t, float64(1), res["id"])
}

func Test_getSentinelByIDHandler(t *testing.T) {
	dbConn := setupTest()
	defer dbConn.Close()

	setupDummyDataSentinelHandler(dbConn)

	w := makeRequest(dbConn, "GET", "/sentinel/2", nil)

	json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 2, res.Data[0].ID)
	assert.Equal(t, 0, len(res.Errors))
}

func Test_getSentinelHandler(t *testing.T) {
	dbConn := setupTest()
	defer dbConn.Close()

	setupDummyDataSentinelHandler(dbConn)

	w := makeRequest(dbConn, "GET", "/sentinel", nil)
	json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 2, len(res.Data))
	assert.Equal(t, 0, len(res.Errors))
}

func Test_updateSentinelHandler(t *testing.T) {
	dbConn := setupTest()
	defer dbConn.Close()

	setupDummyDataSentinelHandler(dbConn)

	updateSentinel := model.Sentinel{
		Name: "update-name",
	}

	var res map[string]any

	w := makeRequest(dbConn, "PATCH", "/sentinel/2", updateSentinel)
	json.Unmarshal(w.Body.Bytes(), &res)

	fmt.Println(res)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, float64(1), res["updated_row"])
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
