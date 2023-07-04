package main

import (
	"context"
	"fmt"

	redis "github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	sentinel := redis.NewSentinelClient(&redis.Options{
		Addr: ":26379",
	})

	addr, err := sentinel.GetMasterAddrByName(ctx, "mymaster").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println(addr)
}
