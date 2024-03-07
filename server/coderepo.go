/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package server provides functionality for setting up and configuring a server for handling code repo operations.
package server

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/coderepo/app"
	"github.com/openmerlin/merlin-server/coderepo/controller"
	"github.com/openmerlin/merlin-server/coderepo/infrastructure/branchclientadapter"
	"github.com/openmerlin/merlin-server/coderepo/infrastructure/branchrepositoryadapter"
	"github.com/openmerlin/merlin-server/coderepo/infrastructure/coderepoadapter"
	"github.com/openmerlin/merlin-server/coderepo/infrastructure/coderepofileadapter"
	"github.com/openmerlin/merlin-server/coderepo/infrastructure/resourceadapterimpl"
	"github.com/openmerlin/merlin-server/common/infrastructure/gitea"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/config"
	"github.com/openmerlin/merlin-server/models/infrastructure/modelrepositoryadapter"
	"github.com/openmerlin/merlin-server/space/infrastructure/spacerepositoryadapter"
)

func initCodeRepo(cfg *config.Config, services *allServices) error {
	err := branchrepositoryadapter.Init(postgresql.DB(), &cfg.CodeRepo.Tables)
	if err != nil {
		return err
	}

	services.codeRepoApp = app.NewCodeRepoAppService(
		coderepoadapter.NewRepoAdapter(gitea.Client(), services.userApp),
	)

	services.codeRepoFileApp = app.NewCodeRepoFileAppService(
		coderepofileadapter.NewCodeRepoFileAdapter(gitea.Client()),
	)

	return nil
}

func initResource(services *allServices) {
	services.resourceApp = app.NewResourceAppService(
		resourceadapterimpl.NewResourceAdapterImpl(
			modelrepositoryadapter.ModelAdapter(),
			spacerepositoryadapter.SpaceAdapter(),
		))
}

func setRouterOfCodeRepoFile(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouterForCodeRepoFileController(
		rg,
		services.codeRepoFileApp,
		services.userMiddleWare,
	)
}

func setRouterOfCodeRepo(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouterForCodeRepoController(
		rg,
		services.resourceApp,
		services.userMiddleWare,
		services.rateLimiterMiddleWare,
	)
}

func setRouterOfBranchRestful(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouteForBranchRestfulController(
		rg,
		app.NewBranchAppService(
			services.permissionApp,
			branchrepositoryadapter.BranchAdapter(),
			resourceadapterimpl.NewResourceAdapterImpl(
				modelrepositoryadapter.ModelAdapter(),
				spacerepositoryadapter.SpaceAdapter(),
			),
			branchclientadapter.NewBranchClientAdapter(gitea.Client()),
		),
		services.userMiddleWare,
		services.operationLog,
		services.rateLimiterMiddleWare,
	)
}

func setRouterOfCodeRepoPermissionInternal(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouteForCodeRepoPermissionInternalController(
		rg,
		services.permissionApp,
		resourceadapterimpl.NewResourceAdapterImpl(
			modelrepositoryadapter.ModelAdapter(),
			spacerepositoryadapter.SpaceAdapter(),
		),
		services.userMiddleWare,
	)
}
