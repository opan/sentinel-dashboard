package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sentinel-dashboard/db"

	"context"

	redis "github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type handler struct {
	DB        db.DB
	GinRouter *gin.Engine
}

type Handler interface {
	Router()
	registerSentinelHandler() gin.HandlerFunc
}

func (h *handler) Router() {
	h.GinRouter.GET("/sentinels", connectSentinel)
	h.GinRouter.POST("/sentinel/register", h.registerSentinelHandler())
}

func (h *handler) Start() {
	h.GinRouter.Run("localhost:2134")
}

func New(dbConn db.DB) handler {
	router := gin.Default()
	h := handler{
		DB:        dbConn,
		GinRouter: router,
	}

	return h
}

func (h *handler) registerSentinelHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	}
}

type Sentinel struct {
	Nodes      []string
	MasterName string
	MasterIP   string
}

func connectSentinel(c *gin.Context) {
	sentinel := redis.NewSentinelClient(&redis.Options{
		Addr: ":26379",
	})

	addr, err := sentinel.GetMasterAddrByName(ctx, "mymaster").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println(addr)
	sentinelRes := Sentinel{
		MasterName: "mymaster",
		MasterIP:   addr[0],
	}

	c.IndentedJSON(http.StatusOK, sentinelRes)
}
