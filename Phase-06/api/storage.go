package main

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var (
	ctx               = context.Background()
	rdb *redis.Client = nil
)

func initRedis() {
	if rdb == nil {
		rdb = redis.NewClient(&redis.Options{
			Addr:     "database:6379",
			Password: "",
			DB:       0,
		})
	}
}

func saveToRedis(key string, value []byte) error {
	if rdb == nil {
		initRedis()
	}

	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
