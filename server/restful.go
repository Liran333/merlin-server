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
	userctl "github.com/openmerlin/merlin-server/user/controller"
)

func setRouterOfRestful(prefix string, engine *gin.Engine, cfg *config.Config, services *allServices) {
	rg := engine.Group(prefix)

	services.securityLog = securitylog.SecurityLog()
	services.userMiddleWare = userctl.RestfulAPI(services.userApp, services.securityLog)
	services.operationLog = operationlog.OperationLog(services.userMiddleWare)
	r, err := ratelimiter.InitRateLimiter(cfg.Redis)
	if err != nil {
		logrus.Fatalf("init ratelimit failed, %s", err)
	}

	services.rateLimiterMiddleWare = r

	// set routers
	setRouterOfOrg(rg, cfg, services)

	setRouterOfUser(rg, cfg, services)

	setRouterOfCodeRepoFile(rg, services)

	setRouterOfModelRestful(rg, services)

	setRouterOfSpaceRestful(rg, services)

	setRouterOfBranchRestful(rg, services)

	api.SwaggerInforest.Title = apiTitle
	api.SwaggerInforest.Version = version
	api.SwaggerInforest.Description = apiDesc
	rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.InstanceName("rest")))
}
