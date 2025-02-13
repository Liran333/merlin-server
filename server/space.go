/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package server provides functionality for setting up and configuring a server for handling code repo operations.
package server

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/common/infrastructure/email"
	"github.com/openmerlin/merlin-server/common/infrastructure/obs"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/common/infrastructure/securestorage"
	"github.com/openmerlin/merlin-server/config"
	modelapp "github.com/openmerlin/merlin-server/models/app"
	"github.com/openmerlin/merlin-server/models/infrastructure/modelrepositoryadapter"
	orgrepoimpl "github.com/openmerlin/merlin-server/organization/infrastructure/repositoryimpl"
	"github.com/openmerlin/merlin-server/space/app"
	"github.com/openmerlin/merlin-server/space/controller"
	emailimpl "github.com/openmerlin/merlin-server/space/infrastructure/emailadapter"
	"github.com/openmerlin/merlin-server/space/infrastructure/messageadapter"
	"github.com/openmerlin/merlin-server/space/infrastructure/obsadapter"
	"github.com/openmerlin/merlin-server/space/infrastructure/securestoragadapter"
	"github.com/openmerlin/merlin-server/space/infrastructure/spacerepositoryadapter"
	"github.com/openmerlin/merlin-server/spaceapp/infrastructure/repositoryadapter"
)

func initSpace(cfg *config.Config, services *allServices) error {
	err := modelrepositoryadapter.Init(postgresql.DB(), &cfg.Model.Tables)
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
		services.npuGatekeeper,
		orgrepoimpl.NewMemberRepo(postgresql.DAO(cfg.Org.Domain.Tables.Member)),
		services.disable,
		services.computilityApp,
		services.spaceappApp,
		services.userApp,
		obsadapter.NewClient(obs.Client()),
		emailimpl.NewEmailImpl(email.GetEmailInst(), cfg.Email.ReportEmail, cfg.Email.RootUrl, cfg.Email.MailTemplate),
	)

	services.modelSpace = app.NewModelSpaceAppService(
		services.permissionApp,
		spacerepositoryadapter.ModelSpaceRelationAdapter(),
		modelrepositoryadapter.ModelAdapter(),
		spacerepositoryadapter.SpaceAdapter(),
		repositoryadapter.AppRepositoryAdapter(),
		modelapp.NewModelInternalAppService(
			modelrepositoryadapter.ModelLabelsAdapter(),
			modelrepositoryadapter.ModelAdapter(),
			modelrepositoryadapter.ModelDeployAdapter(),
		),
	)

	services.spaceVariable = app.NewSpaceVariableService(
		services.permissionApp,
		spacerepositoryadapter.SpaceAdapter(),
		repositoryadapter.AppRepositoryAdapter(),
		spacerepositoryadapter.SpaceVariableAdapter(),
		securestoragadapter.SecureStorageAdapter(securestorage.GetClient(), cfg.Vault.BasePath),
		messageadapter.MessageAdapter(&cfg.Space.Topics),
	)

	services.spaceSecret = app.NewSpaceSecretService(
		services.permissionApp,
		spacerepositoryadapter.SpaceAdapter(),
		repositoryadapter.AppRepositoryAdapter(),
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

func setRouterOfSpaceInternal(rg *gin.RouterGroup, services *allServices, cfg *config.Config) {
	controller.AddRouterForSpaceInternalController(
		rg,
		app.NewSpaceInternalAppService(
			spacerepositoryadapter.SpaceAdapter(),
			messageadapter.MessageAdapter(&cfg.Space.Topics),
			repositoryadapter.AppRepositoryAdapter(),
			spacerepositoryadapter.ModelSpaceRelationAdapter(),
			modelrepositoryadapter.ModelAdapter(),
		),
		services.modelSpace,
		services.userMiddleWare,
	)
}
