package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sentinel-dashboard/model"
)

func (h *handler) registerSentinelHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := h.dbConn.GetConnection()
		tx, err := db.Begin()
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("db.Begin: %w", err))
			return
		}

		stmt, err := tx.Prepare("INSERT INTO sentinels (name, hosts) VALUES (?, ?)")
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("db.Prepare: %w", err))
			return
		}

		defer stmt.Close()

		body := model.Sentinel{}
		if err = ctx.BindJSON(&body); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("BindJSON: %w", err))
			return
		}

		_, err = stmt.Exec(body.Name, body.Hosts)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("stmt.Exec: %w", err))
			return
		}

		err = tx.Commit()
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("tx.Commit: %w", err))
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"msg":    "Sentinel successfully registered",
			"errors": []string{},
		})
	}
}
