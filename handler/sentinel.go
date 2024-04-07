package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sentinel-dashboard/model"
)

func (h *handler) RegisterSentinelHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbx := h.dbConn.GetConnection()
		tx, err := dbx.Beginx()
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("db.Beginx: %w", err))
			return
		}

		stmt, err := tx.Preparex("INSERT INTO sentinels (name, hosts) VALUES ($1, $2)")
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

		res, err := stmt.Exec(body.Name, body.Hosts)
		if err != nil {
			tx.Rollback()
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("stmt.Exec: %w", err))
			return
		}

		lastID, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
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
		var stmtStr string
		var rows *sqlx.Rows

		if id != "" {
			stmtStr = "SELECT * FROM sentinels WHERE id = ?"
		} else {
			stmtStr = "SELECT * FROM sentinels ORDER BY id"
		}

		stmt, err := dbx.Preparex(stmtStr)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("db.Prepare: %w", err))
		}
		defer stmt.Close()

		if id != "" {
			rows, err = stmt.Queryx(id)
			if err != nil {
				ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("stmt.Query: %w", err))
				return
			}
		} else {
			rows, err = stmt.Queryx()
			if err != nil {
				ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("stmt.Query: %w", err))
				return
			}
		}

		defer rows.Close()

		var results []model.Sentinel
		for rows.Next() {
			var r model.Sentinel
			err = rows.StructScan(&r)
			if err != nil {
				ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("rows.Scan: %w", err))
				return
			}

			results = append(results, r)
		}

		err = rows.Err()
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"msg":    "",
			"data":   results,
			"errors": []string{},
		})
	}
}

func (h *handler) UpdateSentinelHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbx := h.dbConn.GetConnection()
		id := ctx.Param("id")

		qs, err := dbx.Preparex("SELECT * FROM sentinels WHERE id = ?")
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("db.Prepare: %w", err))
		}
		defer qs.Close()

		var s model.Sentinel
		err = qs.QueryRowx(id).StructScan(&s)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusOK, gin.H{
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
		if err = ctx.BindJSON(&rb); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("BindJSON: %w", err))
			return
		}

		tx, err := dbx.Beginx()
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("dbx.Begin: %w", err))
			return
		}

		pu, err := dbx.Preparex("UPDATE sentinels SET name = $2, hosts = $3 WHERE id = $1")
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("dbx.Preparex: %w", err))
			return
		}
		defer pu.Close()

		qr, err := pu.Exec(id, rb.Name, rb.Hosts)
		if err != nil {
			tx.Rollback()
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("stmt.Exec: %w", err))
			return
		}

		lastID, err := qr.LastInsertId()
		if err != nil {
			tx.Rollback()
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("LastInsertId: %w", err))
			return
		}

		err = tx.Commit()
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("tx.Exec: %w", err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"msg":    "Record has been successfully updated",
			"id":     lastID,
			"errors": []string{},
		})
	}
}
