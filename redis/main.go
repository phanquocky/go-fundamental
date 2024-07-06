package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	clientRead  *redis.Client
	clientWrite *redis.Client
	*redis.Client
}

// how to implement all Method of redis.Client by using clientRead

func NewRedisClient() *RedisClient {
	clientRead := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
	})

	clientWrite := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
	})

	return &RedisClient{
		clientRead:  clientRead,
		clientWrite: clientWrite,
	}
}

func main() {

	client := NewRedisClient()

	ctx := context.Background()

	err := client.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		fmt.Println("error", err)
	}
}
