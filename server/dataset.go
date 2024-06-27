/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package server provides functionality for setting up and configuring a server for handling datasets operations.
package server

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/common/infrastructure/email"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/config"
	"github.com/openmerlin/merlin-server/datasets/app"
	"github.com/openmerlin/merlin-server/datasets/controller"
	"github.com/openmerlin/merlin-server/datasets/infrastructure/datasetrepositoryadapter"
	"github.com/openmerlin/merlin-server/datasets/infrastructure/emailimpl"
	"github.com/openmerlin/merlin-server/datasets/infrastructure/messageadapter"
	orgrepoimpl "github.com/openmerlin/merlin-server/organization/infrastructure/repositoryimpl"
)

func initDataset(cfg *config.Config, services *allServices) error {
	err := datasetrepositoryadapter.Init(postgresql.DB(), &cfg.Dataset.Tables)
	if err != nil {
		return err
	}

	services.datasetApp = app.NewDatasetAppService(
		services.permissionApp,
		messageadapter.MessageAdapter(&cfg.Dataset.Topics),
		services.codeRepoApp,
		datasetrepositoryadapter.DatasetAdapter(),
		orgrepoimpl.NewMemberRepo(postgresql.DAO(cfg.Org.Domain.Tables.Member)),
		services.disable,
		services.userApp,
		emailimpl.NewEmailImpl(email.GetEmailInst(), cfg.Email.ReportEmail, cfg.Email.RootUrl, cfg.Email.MailTemplate),
	)

	return nil
}

func setRouterOfDatasetWeb(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouteForDatasetWebController(
		rg,
		services.datasetApp,
		services.userMiddleWare,
		services.operationLog,
		services.securityLog,
		services.userApp,
		services.rateLimiterMiddleWare,
		services.privacyCheck,
		services.activityApp,
	)
}

func setRouterOfDatasetRestful(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouteForDatasetRestfulController(
		rg,
		services.datasetApp,
		services.userMiddleWare,
		services.operationLog,
		services.securityLog,
		services.userApp,
		services.rateLimiterMiddleWare,
		services.privacyCheck,
		services.activityApp,
	)
}

func setRouterOfDatasetInternal(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouterForDatasetInternalController(
		rg,
		app.NewDatasetInternalAppService(
			datasetrepositoryadapter.DatasetLabelsAdapter(),
			datasetrepositoryadapter.DatasetAdapter(),
		),
		services.userMiddleWare,
	)
}
