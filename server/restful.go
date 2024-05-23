/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package server provides functionality for setting up and configuring a server for handling code repo operations.
package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/common/controller/middleware/operationlog"
	"github.com/openmerlin/merlin-server/common/controller/middleware/privacycheck"
	"github.com/openmerlin/merlin-server/common/controller/middleware/ratelimiter"
	"github.com/openmerlin/merlin-server/common/controller/middleware/securitylog"
	"github.com/openmerlin/merlin-server/config"
	userctl "github.com/openmerlin/merlin-server/user/controller"
)

func setRouterOfRestful(prefix string, engine *gin.Engine, cfg *config.Config, services *allServices) {
	rg := engine.Group(prefix)

	services.securityLog = securitylog.SecurityLog()
	services.userMiddleWare = userctl.RestfulAPI(services.userApp, services.securityLog)
	services.operationLog = operationlog.OperationLog(services.userMiddleWare)
	services.rateLimiterMiddleWare = ratelimiter.Limiter()
	if services.rateLimiterMiddleWare == nil {
		logrus.Fatalf("init ratelimit failed")
	}

	services.privacyCheck = privacycheck.PrivacyCheck(services.userMiddleWare, services.userApp)

	if cfg.NeedTokenForEachAPI {
		rg.Use(services.userMiddleWare.Read)
	}

	// set routers
	setRouterOfOrg(rg, cfg, services)

	setRouterOfUser(rg, cfg, services)

	setRouterOfModelRestful(rg, services)

	setRouterOfDatasetRestful(rg, services)

	setRouterOfSpaceRestful(rg, services)

	setRouterOfSpaceAppRestful(rg, services)

	setRouterOfBranchRestful(rg, services)

	setRouterOfActivityRestful(rg, services)
}
