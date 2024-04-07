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
	sessionctl "github.com/openmerlin/merlin-server/session/controller"
)

func setRouterOfWeb(prefix string, engine *gin.Engine, cfg *config.Config, services *allServices) {
	rg := engine.Group(prefix)

	services.securityLog = securitylog.SecurityLog()
	services.userMiddleWare = sessionctl.WebAPIMiddleware(services.sessionApp, services.securityLog,
		&cfg.Session.Controller)
	services.tokenMiddleWare = sessionctl.WebAPIMiddleware(services.sessionApp, services.securityLog,
		&cfg.Session.Controller)
	services.operationLog = operationlog.OperationLog(services.userMiddleWare)
	services.rateLimiterMiddleWare = ratelimiter.Limiter()
	if services.rateLimiterMiddleWare == nil {
		logrus.Fatalf("init ratelimit failed")
	}

	services.privacyCheck = privacycheck.PrivacyCheck(services.userMiddleWare, services.userApp)

	// set routers
	setRouterOfOrg(rg, cfg, services)

	setRouterOfUser(rg, cfg, services)

	setRouterOfSession(rg, services, &cfg.Session.Controller)

	setRouterOfModelWeb(rg, services)

	setRouterOfSpaceWeb(rg, services)

	setRouterOfSpaceAppWeb(rg, services, cfg)

	setRouterOfCodeRepo(rg, services)

	setRouterOfSearchWeb(rg, services)

	setRouterOfActivityWeb(rg, services)

	setRouterOfOther(rg, &cfg.OtherConfig)
}
