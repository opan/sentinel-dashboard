package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Name  string
	Hosts string
}

func (h *handler) registerSentinelHandler() gin.HandlerFunc {
	var errMsg string
	respCode := http.StatusOK

	return func(ctx *gin.Context) {
		db := h.dbConn.GetConnection()
		tx, err := db.Begin()
		if err != nil {
			errMsg = err.Error()
		}

		stmt, err := tx.Prepare("INSERT INTO sentinels (name, hosts) VALUES (?, ?)")
		if err != nil {
			errMsg = err.Error()
		}
		defer stmt.Close()

		body := Body{}

		if err = ctx.BindJSON(&body); err != nil {
			errMsg = fmt.Errorf("BindJSON %+v", err).Error()
		}

		_, err = stmt.Exec(body.Name, body.Hosts)
		if err != nil {
			errMsg = err.Error()
		}

		err = tx.Commit()
		if err != nil {
			errMsg = err.Error()
		}

		ctx.JSON(respCode, gin.H{
			"errMsg": errMsg,
		})
	}
}
