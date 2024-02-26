/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

const expireTime = 3

type dbRedis struct {
	Expiration time.Duration
}

// NewDBRedis creates a new dbRedis instance with the given expiration time.
func NewDBRedis(expiration int) dbRedis {
	return dbRedis{Expiration: time.Duration(expiration)}
}

// Create sets the value for the given key with the specified expiration time.
func (r dbRedis) Create(
	ctx context.Context, key string, value interface{},
) *redis.StatusCmd {
	return client.Set(ctx, key, value, r.Expiration*time.Second)
}

// Get retrieves the value for the given key.
func (r dbRedis) Get(
	ctx context.Context, key string,
) *redis.StringCmd {
	return client.Get(ctx, key)
}

// Delete deletes the key from the database.
func (r dbRedis) Delete(
	ctx context.Context, key string,
) *redis.IntCmd {
	return client.Del(ctx, key)
}

// Expire sets the expiration time for the given key.
func (r dbRedis) Expire(
	ctx context.Context, key string, expire time.Duration,
) *redis.BoolCmd {
	return client.Expire(ctx, key, expireTime*time.Second)
}

// DB returns the Redis client instance.
func DB() *redis.Client {
	return client
}
