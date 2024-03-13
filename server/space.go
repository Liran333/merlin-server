/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package server

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/config"
	"github.com/openmerlin/merlin-server/space/app"
	"github.com/openmerlin/merlin-server/space/controller"
	"github.com/openmerlin/merlin-server/space/infrastructure/messageadapter"
	"github.com/openmerlin/merlin-server/space/infrastructure/spacerepositoryadapter"
)

func initSpace(cfg *config.Config, services *allServices) error {
	err := spacerepositoryadapter.Init(postgresql.DB(), &cfg.Space.Tables)
	if err != nil {
		return err
	}

	services.spaceApp = app.NewSpaceAppService(
		services.permissionApp,
		messageadapter.MessageAdapter(&cfg.Space.Topics),
		services.codeRepoApp,
		spacerepositoryadapter.SpaceAdapter(),
	)

	return nil
}

func setRouterOfSpaceWeb(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouteForSpaceWebController(
		rg,
		services.spaceApp,
		services.userMiddleWare,
		services.operationLog,
		services.securityLog,
		services.rateLimiterMiddleWare,
		services.userApp,
		services.privacyCheck,
	)
}

func setRouterOfSpaceRestful(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouteForSpaceRestfulController(
		rg,
		services.spaceApp,
		services.userMiddleWare,
		services.operationLog,
		services.securityLog,
		services.rateLimiterMiddleWare,
		services.userApp,
		services.privacyCheck,
	)
}

func setRouterOfSpaceInternal(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouterForSpaceInternalController(
		rg,
		app.NewSpaceInternalAppService(
			spacerepositoryadapter.SpaceAdapter(),
		),
		services.userMiddleWare,
	)
}
