package server

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/config"

	"github.com/openmerlin/merlin-server/space-app/app"
	"github.com/openmerlin/merlin-server/space-app/controller"
	"github.com/openmerlin/merlin-server/space-app/infrastructure/messageadapter"
	"github.com/openmerlin/merlin-server/space-app/infrastructure/repositoryadapter"
	"github.com/openmerlin/merlin-server/space/infrastructure/spacerepositoryadapter"
)

func initSpaceApp(cfg *config.Config, services *allServices) error {
	return repositoryadapter.Init(postgresql.DB(), &cfg.SpaceApp.Tables)
}

func setRouterOfSpaceAppWeb(rg *gin.RouterGroup, services *allServices) {
	s := app.NewSpaceappAppService(
		repositoryadapter.AppRepositoryAdapter(),
		spacerepositoryadapter.SpaceAdapter(),
		services.permission,
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
