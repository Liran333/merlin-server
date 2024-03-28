/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package server

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/common/controller/middleware/internalservice"
	"github.com/openmerlin/merlin-server/common/controller/middleware/securitylog"
	"github.com/openmerlin/merlin-server/config"
)

func setRouterOfInternal(prefix string, engine *gin.Engine, cfg *config.Config, services *allServices) {
	rg := engine.Group(prefix)

	services.securityLog = securitylog.SecurityLog()
	services.userMiddleWare = internalservice.NewAPIMiddleware(services.securityLog)

	// set routers
	setRouterOfSessionInternal(rg, services)

	setInternalRouterOfUser(rg, cfg, services)

	setRouterOfSpaceInternal(rg, services)

	setRouterOfModelInternal(rg, services)

	setRouterOfActivityInternal(rg, services)

	setRouterOfSpaceAppInternal(rg, services, cfg)

	setRouterOfCodeRepoPermissionInternal(rg, services)

	rg.GET("/heartbeat", func(*gin.Context) {})
}
