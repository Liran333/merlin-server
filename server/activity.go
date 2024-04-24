/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package server provides functionality for setting up and configuring a server for handling code repo operations.
package server

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/activity/app"
	"github.com/openmerlin/merlin-server/activity/controller"
	"github.com/openmerlin/merlin-server/activity/insfrastructure/activityrepositoryadapter"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/config"
)

func initActivity(cfg *config.Config, services *allServices) error {
	err := activityrepositoryadapter.Init(postgresql.DB(), &cfg.Activity.Tables)
	if err != nil {
		return err
	}

	services.activityApp = app.NewActivityAppService(
		services.permissionApp,
		services.codeRepoApp,
		activityrepositoryadapter.ActivityAdapter(),
	)

	return nil
}

func setRouterOfActivityWeb(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouteForActivityWebController(
		rg,
		services.activityApp,
		services.userMiddleWare,
		services.orgApp,
		services.userApp,
		services.modelApp,
		services.spaceApp,
		services.rateLimiterMiddleWare,
		services.operationLog,
	)
}

func setRouterOfActivityRestful(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouteForActivityRestfulController(
		rg,
		services.activityApp,
		services.userMiddleWare,
		services.orgApp,
		services.userApp,
		services.modelApp,
		services.spaceApp,
		services.operationLog,
	)
}

func setRouterOfActivityInternal(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouterForActivityInternalController(
		rg,
		app.NewActivityInternalAppService(
			activityrepositoryadapter.ActivityAdapter(),
		),
		services.userMiddleWare,
	)
}
