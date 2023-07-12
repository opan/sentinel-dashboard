package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) registerSentinelHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	}
}

func (h *handler) listMasterHandler() gin.HandlerFunc {
	masters, err := h.Sentinel.Masters(redisCtx).Result()

	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"err":     err,
			"masters": masters,
		})
	}
}
