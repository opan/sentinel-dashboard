package handler

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sentinel-manager/model"
)

func (h *handler) ClusterInfoHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbx := h.dbConn.GetConnection()
		id := ctx.Param("id")

		var s model.Sentinel

		err := dbx.Get(&s, "SELECT * FROM sentinels WHERE id = ?", id)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"msg":    fmt.Sprintf("No record found with ID: %s", id),
				"data":   nil,
				"errors": []string{},
			})
			return
		}

		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("error get record: %w", err))
			return
		}

		sentinel := redis.NewSentinelClient(&redis.Options{})
		defer sentinel.Close()

		// set 5s timeout when executing
		ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		ping, err := sentinel.Ping(ctxTimeout).Result()
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed ping: %w", err))
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data":   ping,
			"errors": []string{},
		})
	}
}
