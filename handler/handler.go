package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sentinel-dashboard/db"
)

// var redisCtx = context.Background()

type handler struct {
	dbConn    db.DB
	GinRouter *gin.Engine
}

type Handler interface {
	Router()
}

func (h *handler) Router() *gin.Engine {
	h.GinRouter.Use(ErrorMiddleware())
	h.GinRouter.POST("/sentinel/register", h.RegisterSentinelHandler())
	h.GinRouter.GET("/sentinel", h.GetSentinelHandler())

	return h.GinRouter
}

func (h *handler) Start() {
	h.GinRouter.Run("localhost:2134")
}

func New(dbConn db.DB) handler {
	router := gin.Default()
	h := handler{
		dbConn:    dbConn,
		GinRouter: router,
	}

	return h
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		var errMsgs []string

		for _, err := range ctx.Errors {
			errMsgs = append(errMsgs, err.Error())
		}

		if len(errMsgs) > 0 {
			ctx.JSON(-1, gin.H{
				"errors": errMsgs,
			})
		}
	}
}
