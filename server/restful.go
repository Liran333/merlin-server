/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package server

import (
	"github.com/gin-gonic/gin"
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
	api.SwaggerInfo.BasePath = prefix

	rg := engine.Group(api.SwaggerInfo.BasePath)

	services.securityLog = securitylog.SecurityLog()
	services.userMiddleWare = userctl.RestfulAPI(services.userApp, services.securityLog)
	services.operationLog = operationlog.OperationLog(services.userMiddleWare)
	services.rateLimiterMiddleWare = ratelimiter.InitRateLimiter(cfg.Redis)

	// set routers
	setRouterOfOrg(rg, cfg, services)

	setRouterOfUser(rg, cfg, services)

	setRouterOfCodeRepoFile(rg, services)

	setRouterOfModelRestful(rg, services)

	setRouterOfSpaceRestful(rg, services)

	setRouterOfBranchRestful(rg, services)

	rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
