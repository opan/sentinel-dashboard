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

type ErrNoHealthySentinel struct {
	Msg string
}

func (e *ErrNoHealthySentinel) Error() string {
	return fmt.Sprintf("No healthy redis sentinel. Hosts: %s", e.Msg)
}

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

		sentinel, err := getSentinel(ctxTimeout, s.Hosts)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, fmt.Errorf("No successfull ping to all available sentinel hosts: %w", err))
			return
		}

		var masters []model.SentinelMaster

		cmd := redis.NewMapStringInterfaceSliceCmd(ctx, "sentinel", "masters")

		err = sentinel.Process(ctx, cmd)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("error when executing command: %w", err))
			return
		}
		defer sentinel.Close()

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

func (h *handler) ClusterRemoveMasterHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbx := h.dbConn.GetConnection()
		id := ctx.Param("id")
		masterName := ctx.Param("master_name")

		if id == "" || masterName == "" {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("one or more required params missing"))
			return
		}

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

		sentinel, err := getSentinel(ctxTimeout, s.Hosts)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, fmt.Errorf("No successfull ping to all available sentinel hosts: %w", err))
			return
		}

		cmd := sentinel.Remove(ctx, masterName)
		r, err := cmd.Result()

		if redis.HasErrorPrefix(err, "No such master with that name") {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if err != nil || r != "OK" {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("error when fetch the cmd result: %w", err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"msg":    fmt.Sprintf("Master %s has been removed", masterName),
			"errors": []string{},
		})
	}
}

func getSentinel(ctx context.Context, hosts string) (*redis.SentinelClient, error) {
	var sentinel *redis.SentinelClient
	var okHost string
	var pingErr error

	// parse sentinel hosts
	// use the first healthy sentinel host found
	// otherwise, raise error if not found
	sh := strings.Split(hosts, ",")
	for _, v := range sh {
		sentinel = redis.NewSentinelClient(&redis.Options{
			Addr: v,
		})

		ping, err := sentinel.Ping(ctx).Result()
		if err == nil && ping == "PONG" {
			okHost = v
			break
		} else {
			sentinel.Close()
		}
	}

	if okHost == "" {
		pingErr = &ErrNoHealthySentinel{Msg: hosts}
	}

	return sentinel, pingErr
}
