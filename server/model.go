/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package server

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/config"
	"github.com/openmerlin/merlin-server/models/app"
	"github.com/openmerlin/merlin-server/models/controller"
	"github.com/openmerlin/merlin-server/models/infrastructure/messageadapter"
	"github.com/openmerlin/merlin-server/models/infrastructure/modelrepositoryadapter"
)

func initModel(cfg *config.Config, services *allServices) error {
	err := modelrepositoryadapter.Init(postgresql.DB(), &cfg.Model.Tables)
	if err != nil {
		return err
	}

	services.modelApp = app.NewModelAppService(
		services.permissionApp,
		messageadapter.MessageAdapter(&cfg.Model.Topics),
		services.codeRepoApp,
		modelrepositoryadapter.ModelAdapter(),
	)

	return nil
}

func setRouterOfModelWeb(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouteForModelWebController(
		rg,
		services.modelApp,
		services.modelSpace,
		services.userMiddleWare,
		services.operationLog,
		services.securityLog,
		services.userApp,
		services.rateLimiterMiddleWare,
		services.privacyCheck,
		services.activityApp,
	)
}

func setRouterOfModelRestful(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouteForModelRestfulController(
		rg,
		services.modelApp,
		services.userMiddleWare,
		services.operationLog,
		services.securityLog,
		services.userApp,
		services.rateLimiterMiddleWare,
		services.privacyCheck,
		services.activityApp,
	)
}

func setRouterOfModelInternal(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouterForModelInternalController(
		rg,
		app.NewModelInternalAppService(
			modelrepositoryadapter.ModelLabelsAdapter(),
			modelrepositoryadapter.ModelAdapter(),
		),
		services.userMiddleWare,
	)
}
