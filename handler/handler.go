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

func (h *handler) Router() {
	h.GinRouter.POST("/sentinel/register", h.registerSentinelHandler())

	h.GinRouter.Run(":8282")
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
