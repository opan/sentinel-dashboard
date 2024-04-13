package handler

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
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

		// set 5s timeout when executing
		ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		var sentinel *redis.SentinelClient
		var okHost string

		// parse sentinel hosts
		// use the first healthy sentinel host found
		// otherwise, raise error if not found
		sh := strings.Split(s.Hosts, ",")
		for _, v := range sh {
			sentinel = redis.NewSentinelClient(&redis.Options{
				Addr: v,
			})

			ping, err := sentinel.Ping(ctxTimeout).Result()
			if err == nil && ping == "PONG" {
				okHost = v
				break
			} else {
				sentinel.Close()
			}
		}

		if okHost == "" {
			ctx.JSON(http.StatusServiceUnavailable, fmt.Errorf("No successfull ping to all available sentinel hosts: %s", s.Hosts))
			return
		}

		cmd := redis.NewCmd(ctx, "info")

		err = sentinel.Process(ctx, cmd)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("error processing command: %w", err))
			return
		}
		defer sentinel.Close()

		fmt.Println(cmd.Result())

		ctx.JSON(http.StatusOK, gin.H{
			"data":   "",
			"errors": []string{},
		})
	}
}
