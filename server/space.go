/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package server provides functionality for setting up and configuring a server for handling code repo operations.
package server

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/common/infrastructure/securestorage"
	"github.com/openmerlin/merlin-server/config"
	modelapp "github.com/openmerlin/merlin-server/models/app"
	"github.com/openmerlin/merlin-server/models/infrastructure/modelrepositoryadapter"
	"github.com/openmerlin/merlin-server/space/app"
	"github.com/openmerlin/merlin-server/space/controller"
	"github.com/openmerlin/merlin-server/space/infrastructure/messageadapter"
	"github.com/openmerlin/merlin-server/space/infrastructure/securestoragadapter"
	"github.com/openmerlin/merlin-server/space/infrastructure/spacerepositoryadapter"
	"github.com/openmerlin/merlin-server/spaceapp/infrastructure/repositoryadapter"
)

func initSpace(cfg *config.Config, services *allServices) error {
	err := spacerepositoryadapter.Init(postgresql.DB(), &cfg.Space.Tables)
	if err != nil {
		return err
	}

	err = modelrepositoryadapter.Init(postgresql.DB(), &cfg.Model.Tables)
	if err != nil {
		return err
	}

	services.spaceApp = app.NewSpaceAppService(
		services.permissionApp,
		messageadapter.MessageAdapter(&cfg.Space.Topics),
		services.codeRepoApp,
		repositoryadapter.AppRepositoryAdapter(),
		spacerepositoryadapter.SpaceVariableAdapter(),
		spacerepositoryadapter.SpaceSecretAdapter(),
		securestoragadapter.SecureStorageAdapter(securestorage.GetClient(), cfg.Vault.BasePath),
		spacerepositoryadapter.SpaceAdapter(),
	)

	services.modelSpace = app.NewModelSpaceAppService(
		services.permissionApp,
		spacerepositoryadapter.ModelSpaceRelationAdapter(),
		modelrepositoryadapter.ModelAdapter(),
		spacerepositoryadapter.SpaceAdapter(),
		modelapp.NewModelInternalAppService(
			modelrepositoryadapter.ModelLabelsAdapter(),
			modelrepositoryadapter.ModelAdapter(),
		),
		app.NewSpaceInternalAppService(spacerepositoryadapter.SpaceAdapter()),
	)

	services.spaceVariable = app.NewSpaceVariableService(
		services.permissionApp,
		spacerepositoryadapter.SpaceAdapter(),
		spacerepositoryadapter.SpaceVariableAdapter(),
		securestoragadapter.SecureStorageAdapter(securestorage.GetClient(), cfg.Vault.BasePath),
		messageadapter.MessageAdapter(&cfg.Space.Topics),
	)

	services.spaceSecret = app.NewSpaceSecretService(
		services.permissionApp,
		spacerepositoryadapter.SpaceAdapter(),
		spacerepositoryadapter.SpaceSecretAdapter(),
		securestoragadapter.SecureStorageAdapter(securestorage.GetClient(), cfg.Vault.BasePath),
		messageadapter.MessageAdapter(&cfg.Space.Topics),
	)

	return nil
}

func setRouterOfSpaceWeb(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouteForSpaceWebController(
		rg,
		services.spaceApp,
		services.modelSpace,
		services.spaceVariable,
		services.spaceSecret,
		services.userMiddleWare,
		services.operationLog,
		services.securityLog,
		services.rateLimiterMiddleWare,
		services.userApp,
		services.privacyCheck,
		services.activityApp,
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
		services.activityApp,
	)
}

func setRouterOfSpaceInternal(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouterForSpaceInternalController(
		rg,
		app.NewSpaceInternalAppService(
			spacerepositoryadapter.SpaceAdapter(),
		),
		services.modelSpace,
		services.userMiddleWare,
	)
}
