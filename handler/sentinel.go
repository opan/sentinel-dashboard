package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sentinel-dashboard/model"
)

func (h *handler) RegisterSentinelHandler() gin.HandlerFunc {
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

func (h *handler) GetSentinelHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := h.dbConn.GetConnection()
		id := ctx.Param("id")
		var stmtStr string
		var rows *sql.Rows

		if id != "" {
			stmtStr = "SELECT * FROM sentinels WHERE id = ?"
		} else {
			stmtStr = "SELECT * FROM sentinels ORDER BY id"
		}

		stmt, err := db.Prepare(stmtStr)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("db.Prepare: %w", err))
		}
		defer stmt.Close()

		if id != "" {
			rows, err = stmt.Query(id)
			if err != nil {
				ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("stmt.Query: %w", err))
				return
			}
		} else {
			rows, err = stmt.Query()
			if err != nil {
				ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("stmt.Query: %w", err))
				return
			}
		}

		defer rows.Close()

		var results []model.Sentinel
		for rows.Next() {
			var r model.Sentinel
			err = rows.Scan(&r.ID, &r.Name, &r.Hosts, &r.CreatedAt)
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
