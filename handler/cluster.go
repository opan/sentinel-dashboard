package handler

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
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
			ctx.JSON(http.StatusInternalServerError, fmt.Errorf("No successfull ping to all available sentinel hosts: %s", s.Hosts))
			return
		}

		var masters []model.SentinelMaster

		cmd := redis.NewMapStringInterfaceSliceCmd(ctx, "sentinel", "masters")

		err = sentinel.Process(ctx, cmd)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("error when executing command: %w", err))
			return
		}

		cr, err := cmd.Result()
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("error when fetching the result: %w", err))
			return
		}

		for _, v := range cr {
			var master model.SentinelMaster

			decodeConfig := &mapstructure.DecoderConfig{
				WeaklyTypedInput: true,
				Result:           &master,
			}

			decoder, err := mapstructure.NewDecoder(decodeConfig)
			if err != nil {
				panic(err)
			}

			decoder.Decode(v)
			masters = append(masters, master)
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data":   masters,
			"errors": []string{},
		})
	}
}
