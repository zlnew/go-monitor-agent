package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type ClientOptions struct {
	Address  string
	Username string
	Password string
	DB       int
}

func Init(ctx context.Context, opts *ClientOptions) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         opts.Address,
		Username:     opts.Username,
		Password:     opts.Password,
		DB:           opts.DB,
		DialTimeout:  time.Millisecond * 500,
		ReadTimeout:  time.Millisecond * 500,
		WriteTimeout: time.Millisecond * 500,
		MaxRetries:   3,
		PoolSize:     10,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("redis connection failed:", err)
		return nil, fmt.Errorf("redis: failed to connect redis client: %w", err)
	}

	return client, nil
}
