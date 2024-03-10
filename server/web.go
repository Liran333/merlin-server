/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/openmerlin/merlin-server/api"
	"github.com/openmerlin/merlin-server/common/controller/middleware/operationlog"
	"github.com/openmerlin/merlin-server/common/controller/middleware/ratelimiter"
	"github.com/openmerlin/merlin-server/common/controller/middleware/securitylog"
	"github.com/openmerlin/merlin-server/config"
	sessionctl "github.com/openmerlin/merlin-server/session/controller"
)

func setRouterOfWeb(prefix string, engine *gin.Engine, cfg *config.Config, services *allServices) {
	rg := engine.Group(prefix)

	services.securityLog = securitylog.SecurityLog()
	services.userMiddleWare = sessionctl.WebAPIMiddleware(services.sessionApp, services.securityLog)
	services.operationLog = operationlog.OperationLog(services.userMiddleWare)
	services.rateLimiterMiddleWare = ratelimiter.Limiter()
	if services.rateLimiterMiddleWare == nil {
		logrus.Fatalf("init ratelimit failed")
	}

	// set routers
	setRouterOfOrg(rg, cfg, services)

	setRouterOfUser(rg, cfg, services)

	setRouterOfSession(rg, services)

	setRouterOfModelWeb(rg, services)

	setRouterOfSpaceWeb(rg, services)

	setRouterOfSpaceAppWeb(rg, services)

	setRouterOfCodeRepo(rg, services)

	api.SwaggerInfoweb.Title = apiTitle
	api.SwaggerInfoweb.Version = version
	api.SwaggerInfoweb.Description = apiDesc
	rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.InstanceName("web")))
}
