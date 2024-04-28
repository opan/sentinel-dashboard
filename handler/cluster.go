package handler

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/mitchellh/mapstructure"
	"github.com/redis/go-redis/v9"
	"github.com/sentinel-manager/model"
)

type sentinelListOfMasters map[string][]model.SentinelMaster

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

		sh := strings.Split(s.Hosts, ",")

		err = sentinelHealthCheck(ctxTimeout, sh)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, fmt.Errorf("some error occured during ping checking: %w", err))
			return
		}

		sm, err := sentinelGetMasters(ctxTimeout, sh)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, fmt.Errorf("some error occured during fetching masters: %w", err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data":   sm,
			"errors": []string{},
		})
	}
}

// Handler to sync state stored in db with the live state
// this will remove stale master stored in the db
// the live state must have in balance state first
// meaning all sentinel node have the same total masters monitored
// any custom options set to the master will be reset
func (h *handler) ClusterSyncStateHandler() gin.HandlerFunc {
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
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("some error occured when get record: %w", err))
			return
		}

		// set 5s timeout when executing
		ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		sh := strings.Split(s.Hosts, ",")

		// check if any of sentinel hosts failed at ping, cancel monitor process
		err = sentinelHealthCheck(ctxTimeout, sh)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, fmt.Errorf("some error occured during ping checking: %w", err))
			return
		}

		sentinelMasters, err := sentinelGetMasters(ctxTimeout, sh)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("errors when fetch masters list: %w", err))
			return
		}

		// compare total of the master for each host
		// raise error when there is difference
		mastersCount := make(map[int]int)
		for _, s := range sentinelMasters {
			mastersCount[len(s)] = len(s)
		}
		if len(mastersCount) > 1 {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("masters count is not equal for each sentinel host"))
			return
		}

		var mastersName []string

		for _, m := range sentinelMasters[sh[0]] {
			mastersName = append(mastersName, m.MasterName)
		}

		qi, args, err := sqlx.In("SELECT * FROM sentinel_masters WHERE sentinel_id = ? AND name NOT IN (?)", id, mastersName)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}

		qi = dbx.Rebind(qi)
		rows, err := dbx.QueryxContext(ctxTimeout, qi, args...)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}

		// collect stale masters into a slice
		var staleMasters []string
		for rows.Next() {
			var staleMaster model.SentinelMaster
			err = rows.StructScan(&staleMaster)
			if err != nil {
				break
			}

			staleMasters = append(staleMasters, staleMaster.MasterName)
		}

		// err checking from for-loop
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}

		if len(staleMasters) > 0 {
			txs, err := dbx.Beginx()
			if err != nil {
				ctx.JSON(http.StatusBadRequest, err)
				return
			}
			defer txs.Rollback()

			di, args, err := sqlx.In("DELETE FROM sentinel_masters WHERE sentinel_id = ? AND name IN (?)", id, staleMasters)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, fmt.Errorf("some error occured when constructing query: %w", err))
				return
			}
			di = dbx.Rebind(di)
			_, err = txs.ExecContext(ctxTimeout, di, args...)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, fmt.Errorf("some error occured when clean up stale master: %w", err))
				return
			}

			err = txs.Commit()
			if err != nil {
				ctx.JSON(http.StatusBadRequest, err)
				return
			}
		}

		var syncErr error

		// loop through masters from the first sentinel host
		for _, s := range sentinelMasters[sh[0]] {
			var syncState = func() error {
				tx, err := dbx.Beginx()
				if err != nil {
					return err
				}
				defer tx.Rollback()

				_, err = tx.Exec("DELETE FROM sentinel_masters WHERE sentinel_id = ? AND name = ?", id, s.MasterName)
				if err != nil {
					return err
				}

				_, err = tx.Exec("INSERT INTO sentinel_masters (sentinel_id, name, ip, port, quorum) VALUES (?, ?, ?, ?, ?)",
					id, s.MasterName, s.IP, s.Port, s.Quorum)
				if err != nil {
					return err
				}

				err = tx.Commit()
				if err != nil {
					return err
				}

				return nil
			}

			err = syncState()
			if err != nil {
				syncErr = err
				break
			}
		}

		if syncErr != nil {
			ctx.JSON(http.StatusBadRequest, fmt.Errorf("some error occured while comparing the state: %w", err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data":   sentinelMasters,
			"msg":    "Sentinel cluster state has been successfully synced",
			"errors": []string{},
		})
	}
}

func (h *handler) ClusterAddMasterHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbx := h.dbConn.GetConnection()
		id := ctx.Param("id")

		tx, err := dbx.Beginx()
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		defer tx.Rollback()

		var s model.Sentinel

		err = dbx.Get(&s, "SELECT * FROM sentinels WHERE id = ?", id)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"msg":    fmt.Sprintf("no record found with ID: %s", id),
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

		sh := strings.Split(s.Hosts, ",")

		err = sentinelHealthCheck(ctxTimeout, sh)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, fmt.Errorf("some error occured during ping checking: %w", err))
			return
		}

		var body model.SentinelMaster
		if err = ctx.BindJSON(&body); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// register master for each sentinel
		for _, h := range sh {
			sentinel := redis.NewSentinelClient(&redis.Options{
				Addr: h,
			})

			monCmd := redis.NewStringCmd(ctx, "sentinel", "monitor", body.MasterName, body.IP, body.Port, body.Quorum)

			err = sentinel.Process(ctx, monCmd)
			if err != nil {
				break
			}

			// result should return "OK"
			_, err := monCmd.Result()
			if err != nil {
				break
			}
		}

		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("some error occured when monitor master: %w", err))
			return
		}

		res, err := tx.Exec("INSERT INTO sentinel_masters (sentinel_id, name, ip, port, quorum) VALUES (?, ?, ?, ?, ?)",
			id, body.MasterName, body.IP, body.Port, body.Quorum,
		)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("error when recording master: %w", err))
			return
		}

		lastID, err := res.LastInsertId()
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		err = tx.Commit()
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"msg":    fmt.Sprintf("Master %s has been successfully registered to sentinel", body.MasterName),
			"id":     lastID,
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

		sh := strings.Split(s.Hosts, ",")

		err = sentinelHealthCheck(ctx, sh)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, fmt.Errorf("some error occured during ping checking: %w", err))
			return
		}

		var removeErr error

		for _, h := range sh {
			sentinel := redis.NewSentinelClient(&redis.Options{
				Addr: h,
			})

			cmd := sentinel.Remove(ctxTimeout, masterName)
			r, err := cmd.Result()
			if redis.HasErrorPrefix(err, "No such master with that name") {
				removeErr = fmt.Errorf("no such master with that name %s", masterName)
				break
			}

			if err != nil || r != "OK" {
				removeErr = fmt.Errorf("error when removing master from host %s: %w", h, err)
				break
			}
		}

		if removeErr != nil {
			ctx.AbortWithError(http.StatusBadRequest, removeErr)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"msg":    fmt.Sprintf("Master %s has been removed", masterName),
			"errors": []string{},
		})
	}
}

func sentinelGetMasters(ctx context.Context, hosts []string) (sentinelListOfMasters, error) {
	var cmdErr error

	sentinelMasters := make(sentinelListOfMasters)

	for _, h := range hosts {
		var masters []model.SentinelMaster
		cmd := redis.NewMapStringInterfaceSliceCmd(ctx, "sentinel", "masters")

		sentinel := redis.NewSentinelClient(&redis.Options{Addr: h})

		cmdErr = sentinel.Process(ctx, cmd)
		if cmdErr != nil {
			sentinel.Close()
			break
		}

		cr, cmdErr := cmd.Result()
		if cmdErr != nil {
			sentinel.Close()
			break
		}

		for _, r := range cr {
			var master model.SentinelMaster

			dc := &mapstructure.DecoderConfig{
				WeaklyTypedInput: true,
				Result:           &master,
			}

			decode, err := mapstructure.NewDecoder(dc)
			if err != nil {
				cmdErr = err
				break
			}

			decode.Decode(r)
			masters = append(masters, master)
		}

		if cmdErr != nil {
			sentinel.Close()
			break
		}

		sentinelMasters[h] = masters
		sentinel.Close()
	}

	if cmdErr != nil {
		return sentinelMasters, cmdErr
	}

	return sentinelMasters, nil
}

func sentinelHealthCheck(ctx context.Context, hosts []string) error {
	var pingErr error

	for _, h := range hosts {
		sentinel := redis.NewSentinelClient(&redis.Options{
			Addr: h,
		})

		pong, pingErr := sentinel.Ping(ctx).Result()

		sentinel.Close()

		if pong != "PONG" || pingErr != nil {
			break
		}
	}

	return pingErr
}
