package db

import (
	"context"
	"fmt"
	"log"
	"monitoring/config"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	*redis.Client
}

func NewRedisClient(cfg *config.RedisConfig,db int) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Println("Successfully connected to Redis")
	return &Client{rdb}, nil
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.Client.Set(ctx, key, value, expiration).Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

func (c *Client) Del(ctx context.Context, key string) error {
	return c.Client.Del(ctx, key).Err()
}

func (c *Client) Exists(ctx context.Context, key string) (int64, error) {
	return c.Client.Exists(ctx, key).Result()
}
