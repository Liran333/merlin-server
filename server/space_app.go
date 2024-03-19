/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package server

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/config"
	"github.com/openmerlin/merlin-server/space/infrastructure/spacerepositoryadapter"
	"github.com/openmerlin/merlin-server/spaceapp/app"
	"github.com/openmerlin/merlin-server/spaceapp/controller"
	"github.com/openmerlin/merlin-server/spaceapp/infrastructure/messageadapter"
	"github.com/openmerlin/merlin-server/spaceapp/infrastructure/repositoryadapter"
	"github.com/openmerlin/merlin-server/spaceapp/infrastructure/sseadapter"
)

func initSpaceApp(cfg *config.Config, services *allServices) error {
	return repositoryadapter.Init(postgresql.DB(), &cfg.SpaceApp.Tables)
}

func setRouterOfSpaceAppWeb(rg *gin.RouterGroup, services *allServices, cfg *config.Config) {
	s := app.NewSpaceappAppService(
		messageadapter.MessageAdapter(&cfg.SpaceApp.Topics),
		repositoryadapter.AppRepositoryAdapter(),
		spacerepositoryadapter.SpaceAdapter(),
		services.permissionApp,
		sseadapter.StreamSentAdapter(),
	)

	controller.AddRouterForSpaceappWebController(
		rg,
		s,
		services.userMiddleWare,
		services.tokenMiddleWare,
		services.rateLimiterMiddleWare,
	)
}

func setRouterOfSpaceAppRestful(rg *gin.RouterGroup, services *allServices, cfg *config.Config) {
	s := app.NewSpaceappAppService(
		messageadapter.MessageAdapter(&cfg.SpaceApp.Topics),
		repositoryadapter.AppRepositoryAdapter(),
		spacerepositoryadapter.SpaceAdapter(),
		services.permissionApp,
		sseadapter.StreamSentAdapter(),
	)

	controller.AddRouterForSpaceappRestfulController(
		rg,
		s,
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
	)

	controller.AddRouteForSpaceappInternalController(
		rg, s, services.userMiddleWare,
	)
}
