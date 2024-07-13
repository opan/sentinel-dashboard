package handler

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sentinel-manager/model"
)

func (h *handler) CreateSentinelHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbx := h.dbConn.GetConnection()
		tx, err := dbx.Beginx()
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("db.Beginx: %w", err))
			return
		}
		defer tx.Rollback()

		body := model.Sentinel{}
		if err = ctx.BindJSON(&body); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("BindJSON: %w", err))
			return
		}

		res, err := tx.Exec("INSERT INTO sentinels (name, hosts) VALUES (?, ?)", body.Name, body.Hosts)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("stmt.Exec: %w", err))
			return
		}

		lastID, err := res.LastInsertId()
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("res.LastInsertId: %w", err))
			return
		}

		err = tx.Commit()
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("tx.Commit: %w", err))
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"msg":    "Sentinel successfully registered",
			"id":     lastID,
			"errors": []string{},
		})
	}
}

func (h *handler) GetSentinelHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbx := h.dbConn.GetConnection()
		id := ctx.Param("id")
		var qs string
		var results []model.Sentinel

		if id != "" {
			qs = "SELECT * FROM sentinels WHERE id = ?"
		} else {
			qs = "SELECT * FROM sentinels ORDER BY id"
		}

		err := dbx.Select(&results, qs, id)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("query select: %w", err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data":   results,
			"errors": []string{},
		})
	}
}

func (h *handler) UpdateSentinelHandler() gin.HandlerFunc {
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
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("rows.StructScan: %w", err))
			return
		}

		rb := model.Sentinel{}
		if err := ctx.BindJSON(&rb); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("BindJSON: %w", err))
			return
		}

		tx, err := dbx.Beginx()
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("dbx.Begin: %w", err))
			return
		}
		defer tx.Rollback()

		r, err := tx.Exec("UPDATE sentinels SET name = ?, hosts = ? WHERE id = ?", rb.Name, rb.Hosts, id)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("stmt.Exec: %w", err))
			return
		}

		ur, err := r.RowsAffected()
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("fetch affetcted row: %w", err))
			return
		}

		err = tx.Commit()
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("tx.Exec: %w", err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"msg":         "Record has been successfully updated",
			"updated_row": ur,
			"errors":      []string{},
		})
	}
}

func (h *handler) RemoveSentinelHandler() gin.HandlerFunc {
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
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("errors get records: %w", err))
			return
		}

		// set 5s timeout when executing
		ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		tx, err := dbx.Beginx()
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("errors begin trx: %w", err))
			return
		}

		defer tx.Rollback()
		_, err = tx.ExecContext(ctxTimeout, "DELETE FROM sentinel_masters WHERE sentinel_id = ?", id)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("errors removing sentinel masters: %w", err))
			return
		}

		_, err = tx.ExecContext(ctxTimeout, "DELETE FROM sentinels WHERE id = ?", id)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("errors removing sentinels: %w", err))
			return
		}

		err = tx.Commit()
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("errors comitting trx: %w", err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"msg":    fmt.Sprintln("Sentinel servers has been removed"),
			"errors": []string{},
		})
	}
}
