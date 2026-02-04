package config

import (
	"context"
	"go-shopping-cart/internal/utils"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addrs    string
	Username string
	Password string
	DB       int
}

func NewRedisClient() *redis.Client {
	cfg := RedisConfig{
		Addrs:    utils.GetEnv("REDIS_ADDR", "localhost:6379"),
		Username: utils.GetEnv("REDIS_USER", "default"),
		Password: utils.GetEnv("REDIS_PASSWORD", "123456"),
		DB:       utils.GetIntEnv("REDIS_DB", 0),
	}

	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addrs,
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     20,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")

	return client
}
