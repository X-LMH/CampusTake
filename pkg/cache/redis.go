package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr, password string, db, poolSize, minIdleConns int) *redis.Client {
	if poolSize <= 0 {
		poolSize = 100
	}
	if minIdleConns <= 0 {
		minIdleConns = poolSize / 5
		if minIdleConns < 5 {
			minIdleConns = 5
		}
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  3 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
	})

	// 启动时探活（可选但推荐）
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = rdb.Ping(ctx).Err() // 你也可以改成失败直接 panic
	return rdb
}
