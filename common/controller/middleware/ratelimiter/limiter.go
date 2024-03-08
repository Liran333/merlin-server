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
	userIdParsed          = "user_id"
	defaultClientPoolSize = 10
	defaultIdleTimeOutNum = 30
	RequestNumPerSec      = 10
	BurstNumPerSec        = 10
)

var (
	overLimitExec = allerror.NewOverLimit(allerror.ErrorRateLimitOver, "too many requests")
)

// InitRateLimiter creates a new instance of the rateLimiter struct.
func InitRateLimiter(cfg redislib.Config) (*rateLimiter, error) {
	// Initialize a redis client using go-redis
	client := &redis.Client{}
	if cfg.DBCert != "" {
		ca, err := ioutil.ReadFile(cfg.DBCert)
		if err != nil {
			return nil, fmt.Errorf("read cert failed")
		}

		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(ca) {
			return nil, fmt.Errorf("new pool failed")
		}

		tlsConfig := &tls.Config{
			InsecureSkipVerify: true, // #nosec G402
			RootCAs:            pool,
		}

		client = redis.NewClient(&redis.Options{
			PoolSize:    defaultClientPoolSize, // default
			IdleTimeout: defaultIdleTimeOutNum * time.Second,
			DB:          cfg.DB,
			Addr:        cfg.Address,
			Password:    cfg.Password,
			TLSConfig:   tlsConfig,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			PoolSize:    defaultClientPoolSize, // default
			IdleTimeout: defaultIdleTimeOutNum * time.Second,
			DB:          cfg.DB,
			Addr:        cfg.Address,
			Password:    cfg.Password,
		})
	}
	// Setup store
	store, err := goredisstore.NewCtx(client, "api-rate-limit:")
	if err != nil {
		logrus.Infof("get new redis client err:%s", err)
		return nil, fmt.Errorf("init goredisstore failed, %s", err)
	}

	requestNum := RequestNumPerSec
	if config.RequestNum > 0 {
		requestNum = config.RequestNum
	}
	burstNum := BurstNumPerSec
	if config.BurstNum > 0 {
		burstNum = config.BurstNum
	}
	// Setup quota
	quota := throttled.RateQuota{
		MaxRate:  throttled.PerSec(requestNum),
		MaxBurst: burstNum,
	}
	rateLimitCtx, err := throttled.NewGCRARateLimiterCtx(store, quota)
	if err != nil {
		logrus.Errorf("get new rate store err:%s", err)
		return nil, fmt.Errorf("init NewGCRARateLimiterCtx failed, %s", err)
	}

	httpRateLimiter := &throttled.HTTPRateLimiterCtx{
		RateLimiter: rateLimitCtx,
	}
	return &rateLimiter{limitCli: httpRateLimiter}, nil
}

type rateLimiter struct {
	limitCli *throttled.HTTPRateLimiterCtx
}

func (rl *rateLimiter) CheckLimit(ctx *gin.Context) {
	v, ok := ctx.Get(userIdParsed)
	if !ok {
		ctx.Next()
		return
	}

	key := fmt.Sprintf("%v", v)
	logrus.Infof("user %v", v)
	logrus.Infof("rl %v", rl)
	logrus.Infof("limitcli %v", rl.limitCli)
	logrus.Infof("ratelimiter %v", rl.limitCli.RateLimiter)
	logrus.Infof("ctx %v", ctx)
	logrus.Infof("ctx request %v", ctx.Request)

	limited, _, err := rl.limitCli.RateLimiter.RateLimitCtx(ctx.Request.Context(), key, 1)
	if err != nil {
		logrus.Errorf("check limit is err:%s", err)
		ctx.Abort()
		return
	}

	if limited {
		commonctl.SendError(ctx, overLimitExec)
		ctx.Abort()
	} else {
		ctx.Next()
	}
}
