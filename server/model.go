/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package server provides functionality for setting up and configuring a server for handling code repo operations.
package server

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/common/infrastructure/email"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/config"
	"github.com/openmerlin/merlin-server/models/app"
	"github.com/openmerlin/merlin-server/models/controller"
	"github.com/openmerlin/merlin-server/models/infrastructure/emailimpl"
	"github.com/openmerlin/merlin-server/models/infrastructure/messageadapter"
	"github.com/openmerlin/merlin-server/models/infrastructure/modelrepositoryadapter"
	orgrepoimpl "github.com/openmerlin/merlin-server/organization/infrastructure/repositoryimpl"
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
		orgrepoimpl.NewMemberRepo(postgresql.DAO(cfg.Org.Domain.Tables.Member)),
		services.disable,
		services.userApp,
		emailimpl.NewEmailImpl(email.GetEmailInst(), cfg.Email.ReportEmail, cfg.Email.RootUrl, cfg.Email.MailTemplate),
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
		services.modelSpace,
		services.userMiddleWare,
	)
}
