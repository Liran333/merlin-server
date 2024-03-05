/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package ratelimiter provides functionality for logging operation-related information.
package ratelimiter

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	redislib "github.com/opensourceways/redis-lib"
	"github.com/sirupsen/logrus"
	"github.com/throttled/throttled/v2"
	"github.com/throttled/throttled/v2/store/goredisstore"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
)

const (
	userIdParsed = "user_id"
)

var (
	overLimitExec = allerror.NewOverLimit(allerror.ErrorRateLimitOver, "request is over limit")
)

// InitRateLimiter creates a new instance of the operationLog struct.
func InitRateLimiter(cfg redislib.Config) *rateLimiter {
	// Initialize a redis client using go-redis
	client := &redis.Client{}
	if cfg.DBCert != "" {
		ca, err := ioutil.ReadFile(cfg.DBCert)
		if err != nil {
			return nil
		}

		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(ca) {
			return nil
		}

		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			RootCAs:            pool,
		}

		client = redis.NewClient(&redis.Options{
			PoolSize:    10, // default
			IdleTimeout: 30 * time.Second,
			DB:          cfg.DB,
			Addr:        cfg.Address,
			Password:    cfg.Password,
			TLSConfig:   tlsConfig,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			PoolSize:    10, // default
			IdleTimeout: 30 * time.Second,
			DB:          cfg.DB,
			Addr:        cfg.Address,
			Password:    cfg.Password,
		})
	}
	// Setup store
	store, err := goredisstore.NewCtx(client, "api-rate-limit:")
	if err != nil {
		// TODO
		logrus.Infof("new redis client err:%s", err)
		return nil
	}
	// Setup quota
	quota := throttled.RateQuota{
		MaxRate:  throttled.PerMin(1),
		MaxBurst: 1,
	}
	rateLimterCtx, err := throttled.NewGCRARateLimiterCtx(store, quota)
	if err != nil {
		// TODO
		logrus.Infof("new rate store err:%s", err)
		return nil
	}

	httpRateLimiter := &throttled.HTTPRateLimiterCtx{
		RateLimiter: rateLimterCtx,
	}
	return &rateLimiter{limitCli: httpRateLimiter}
}

type rateLimiter struct {
	limitCli *throttled.HTTPRateLimiterCtx
}

func (rl *rateLimiter) CheckLimit(ctx *gin.Context) {
	v, ok := ctx.Get(userIdParsed)
	logrus.Infof("get user is :%s, ok :%v", v, ok)
	if !ok {
		logrus.Infof("is checkout ok :%s", ok)
		ctx.Next()
		return
	}
	key := fmt.Sprintf("%v", v)

	limited, _, err := rl.limitCli.RateLimiter.RateLimitCtx(ctx.Request.Context(), key, 1)

	if err != nil {
		// TODO
		logrus.Infof("rate limit key:%s, err:%s", key, err)
		return
	}

	if limited {
		logrus.Infof("rate limit key:%s", key)
		commonctl.SendError(ctx, overLimitExec)
		logrus.Infof("rate limit limit:%s", limited)
		ctx.Abort()
	} else {
		ctx.Next()
	}

}
