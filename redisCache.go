package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client
	ctx context.Context
)

func ConnectRedis() {
	ctx = context.Background()

	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis-18578.c90.us-east-1-3.ec2.redns.redis-cloud.com:18578",
		Username: "default",
		Password: "N2iePQAw4s0jSKhcwP2NIpV79F6IvFei",
		DB:       0,
	})
}

func SetRedisData(key string, data interface{}) {
	rdb.Set(ctx, key, data, 0)
}

func GetRedisData(key string) string {
	res, _ := rdb.Get(ctx, key).Result()
	return res
}
