package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/joho/godotenv"
)

var (
	rdb *redis.Client
	ctx context.Context
)

func ConnectRedis() {

	err := godotenv.Load()
	if err != nil{
		fmt.Println("error loading env file for redis")
	}

	ctx = context.Background()

	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Username: "default",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}

func SetRedisData(key string, data interface{}) {
	rdb.Set(ctx, key, data, 6 * time.Hour)
}

func GetRedisData(key string) string {
	res, _ := rdb.Get(ctx, key).Result()
	return res
}
