package server

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/computility/app"
	"github.com/openmerlin/merlin-server/computility/controller"
	"github.com/openmerlin/merlin-server/computility/infrastructure/messageadapter"
	"github.com/openmerlin/merlin-server/computility/infrastructure/repositoryadapter"
	"github.com/openmerlin/merlin-server/config"
)

func initComputilityApp(cfg *config.Config, services *allServices) error {
	err := repositoryadapter.Init(postgresql.DB(), &cfg.Computility.Tables)
	services.computilityApp = app.NewComputilityInternalAppService(
		repositoryadapter.ComputilityOrgAdapter(),
		repositoryadapter.ComputilityDetailAdapter(),
		repositoryadapter.ComputilityAccountAdapter(),
		repositoryadapter.ComputilityAccountRecordAdapter(),
		messageadapter.MessageAdapter(&cfg.Computility.Topics),
	)

	return err
}

func setRouterOfComputilityAppInternal(rg *gin.RouterGroup, services *allServices, cfg *config.Config) {
	controller.AddRouterForComputilityInternalController(
		rg,
		services.computilityApp,
		services.userMiddleWare,
	)
}

func setRouterOfComputilityAppWeb(rg *gin.RouterGroup, services *allServices) {
	s := app.NewComputilityAppService(
		repositoryadapter.ComputilityOrgAdapter(),
		repositoryadapter.ComputilityDetailAdapter(),
		repositoryadapter.ComputilityAccountAdapter(),
	)

	controller.AddRouterForComputilityWebController(
		rg,
		s,
		services.userMiddleWare,
		services.operationLog,
	)
}
