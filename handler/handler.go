package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"context"

	redis "github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func Router() {
	router := gin.Default()
	router.GET("/sentinels", connectSentinel)
	router.Run("localhost:2134")
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
