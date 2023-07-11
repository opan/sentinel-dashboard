package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sentinel-dashboard/db"

	"context"

	redis "github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type handler struct {
	DB        db.DB
	GinRouter *gin.Engine
	Sentinel  *redis.SentinelClient
}

type Handler interface {
	Router()
	registerSentinelHandler() gin.HandlerFunc
}

func (h *handler) Router() {
	h.GinRouter.GET("/sentinels", h.listMasterHandler())
	h.GinRouter.POST("/sentinel/register", h.registerSentinelHandler())
}

func (h *handler) Start() {
	h.GinRouter.Run("localhost:2134")
}

func New(dbConn db.DB, sentinel *redis.SentinelClient) handler {
	router := gin.Default()
	h := handler{
		DB:        dbConn,
		GinRouter: router,
		Sentinel:  sentinel,
	}

	return h
}
