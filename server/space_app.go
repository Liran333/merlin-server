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
)

func initSpaceApp(cfg *config.Config, services *allServices) error {
	return repositoryadapter.Init(postgresql.DB(), &cfg.SpaceApp.Tables)
}

func setRouterOfSpaceAppWeb(rg *gin.RouterGroup, services *allServices) {
	s := app.NewSpaceappAppService(
		repositoryadapter.AppRepositoryAdapter(),
		spacerepositoryadapter.SpaceAdapter(),
		services.permissionApp,
	)

	controller.AddRouterForSpaceappWebController(
		rg,
		s,
		services.userMiddleWare,
	)
}

func setRouterOfSpaceAppInternal(rg *gin.RouterGroup, services *allServices, cfg *config.Config) {
	s := app.NewSpaceappInternalAppService(
		messageadapter.MessageAdapter(&cfg.SpaceApp.Topics),
		repositoryadapter.AppRepositoryAdapter(),
	)

	controller.AddRouteForSpaceappInternalController(
		rg, s, services.userMiddleWare,
	)
}
