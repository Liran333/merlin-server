/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repositories provides implementations of access repositories using Redis.
package repositories

import (
	"context"
	"time"

	coredis "github.com/openmerlin/merlin-server/common/infrastructure/redis"
	"github.com/openmerlin/merlin-server/infrastructure/redis"
)

// Access represents an access repository interface.
type Access interface {
	Insert(key, value string) error
	Get(key string) (string, error)
	Expire(key string, expire int64) error
}

// NewAccessRepo creates a new access repository with the specified expiration duration.
func NewAccessRepo(expireDuration int) Access {
	return &accessRepo{cli: coredis.NewDBRedis(expireDuration)}
}

type accessRepo struct {
	cli redis.RedisClient
}

// Insert inserts a key-value pair into the access repository.
func (impl *accessRepo) Insert(key, value string) error {
	f := func(ctx context.Context) error {
		cmd := impl.cli.Create(ctx, key, value)
		if cmd.Err() != nil {
			return cmd.Err()
		}

		ok, err := cmd.Result()
		if ok != "ok" {
			return err
		}

		return nil
	}

	return redis.WithContext(f)
}

// Get retrieves the value associated with the specified key from the access repository.
func (impl *accessRepo) Get(key string) (string, error) {
	var value string

	f := func(ctx context.Context) error {
		cmd := impl.cli.Get(ctx, key)
		if cmd.Err() != nil {
			return cmd.Err()
		}

		value = cmd.Val()

		return nil
	}

	if err := redis.WithContext(f); err != nil {
		return "", err
	}

	return value, nil
}

// Expire sets an expiration duration for the specified key in the access repository.
func (impl *accessRepo) Expire(key string, expire int64) error {
	f := func(ctx context.Context) error {
		cmd := impl.cli.Expire(ctx, key, time.Duration(expire))

		return cmd.Err()
	}

	return redis.WithContext(f)
}
