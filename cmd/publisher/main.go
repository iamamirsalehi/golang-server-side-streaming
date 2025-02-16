package main

import (
	"context"
	"fmt"
	golang_server_side_streaming "github.com/server-side-streaming"
	"time"
)

func main() {
	ctx := context.Background()

	redisClient := golang_server_side_streaming.NewRedisClient(ctx)
	channelName := fmt.Sprintf("notification:%s", "123")

	ticker := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down")
		case t := <-ticker.C:
			if cmd := redisClient.Publish(ctx, channelName, fmt.Sprintf("New notification %s", t.String())); cmd.Err() != nil {
				panic(cmd.Err())
			}
		}
	}
}
