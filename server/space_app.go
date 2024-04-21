/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package server provides functionality for setting up and configuring a server for handling code repo operations.
package server

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/config"
	"github.com/openmerlin/merlin-server/space/infrastructure/spacerepositoryadapter"
	"github.com/openmerlin/merlin-server/spaceapp/app"
	"github.com/openmerlin/merlin-server/spaceapp/controller"
	"github.com/openmerlin/merlin-server/spaceapp/infrastructure/messageadapter"
	"github.com/openmerlin/merlin-server/spaceapp/infrastructure/repositoryadapter"
	"github.com/openmerlin/merlin-server/spaceapp/infrastructure/sseadapter"
)

func initSpaceApp(cfg *config.Config, services *allServices) error {

	services.spaceappApp = app.NewSpaceappAppService(
		messageadapter.MessageAdapter(&cfg.SpaceApp.Topics),
		repositoryadapter.AppRepositoryAdapter(),
		spacerepositoryadapter.SpaceAdapter(),
		services.permissionApp,
		sseadapter.StreamSentAdapter(),
		services.compUtilityApp,
	)

	return nil
}

func setRouterOfSpaceAppWeb(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouterForSpaceappWebController(
		rg,
		services.spaceappApp,
		services.userMiddleWare,
		services.tokenMiddleWare,
		services.rateLimiterMiddleWare,
	)
}

func setRouterOfSpaceAppRestful(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouterForSpaceappRestfulController(
		rg,
		services.spaceappApp,
		services.userMiddleWare,
		services.tokenMiddleWare,
		services.rateLimiterMiddleWare,
	)
}

func setRouterOfSpaceAppInternal(rg *gin.RouterGroup, services *allServices, cfg *config.Config) {
	s := app.NewSpaceappInternalAppService(
		messageadapter.MessageAdapter(&cfg.SpaceApp.Topics),
		repositoryadapter.AppRepositoryAdapter(),
		repositoryadapter.BuildLogAdapter(),
		spacerepositoryadapter.SpaceAdapter(),
		services.compUtilityApp,
	)

	controller.AddRouteForSpaceappInternalController(
		rg, s, services.userMiddleWare,
	)
}
