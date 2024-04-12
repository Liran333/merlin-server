/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package server provides functionality for setting up and configuring a server for handling code repo operations.
package server

import (
	"github.com/gin-gonic/gin"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/config"
	orgrepoimpl "github.com/openmerlin/merlin-server/organization/infrastructure/repositoryimpl"

	"github.com/openmerlin/merlin-server/models/infrastructure/modelrepositoryadapter"
	"github.com/openmerlin/merlin-server/search/app"
	"github.com/openmerlin/merlin-server/search/controller"
	"github.com/openmerlin/merlin-server/search/infrastructure/resourceadapterimpl"
	"github.com/openmerlin/merlin-server/space/infrastructure/spacerepositoryadapter"
)

func setRouterOfSearchWeb(rg *gin.RouterGroup, cfg *config.Config, services *allServices) {
	controller.AddRouteForSearchWebController(
		rg,
		app.NewSearchAppService(
			resourceadapterimpl.NewSearchRepositoryAdapter(
				modelrepositoryadapter.ModelAdapter(),
				spacerepositoryadapter.SpaceAdapter(),
				services.userRepo,
				orgrepoimpl.NewMemberRepo(postgresql.DAO(cfg.Org.Tables.Member)),
			),
		),
		services.operationLog,
		services.userMiddleWare,
		services.rateLimiterMiddleWare,
	)
}
