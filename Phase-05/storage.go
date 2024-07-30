package main

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var (
	ctx               = context.Background()
	rdb *redis.Client = nil
)

func initRedis() error {
	if rdb == nil {
		opts, err := redis.ParseURL(redisConnStr)
		if err != nil {
			return err
		}

		rdb = redis.NewClient(opts)
		return rdb.Ping(ctx).Err()
	}

	return nil
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
