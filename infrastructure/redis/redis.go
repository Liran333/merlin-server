/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package redis provides a Redis client interface and utility functions for working with Redis.
package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient is an interface that defines methods for interacting with Redis.
type RedisClient interface {
	Create(context.Context, string, interface{}) *redis.StatusCmd
	Get(context.Context, string) *redis.StringCmd
	Delete(context.Context, string) *redis.IntCmd
	Expire(context.Context, string, time.Duration) *redis.BoolCmd
}

// WithContext is a utility function that executes a function with a context.
func WithContext(f func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	return f(ctx)
}
