/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package server provides functionality for setting up and configuring a server for handling code repo operations.
package server

import (
	activityapp "github.com/openmerlin/merlin-server/activity/app"
	coderepoapp "github.com/openmerlin/merlin-server/coderepo/app"
	commonapp "github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	computilityapp "github.com/openmerlin/merlin-server/computility/app"
	"github.com/openmerlin/merlin-server/config"
	datasetapp "github.com/openmerlin/merlin-server/datasets/app"
	modelapp "github.com/openmerlin/merlin-server/models/app"
	orgapp "github.com/openmerlin/merlin-server/organization/app"
	sessionapp "github.com/openmerlin/merlin-server/session/app"
	spaceapp "github.com/openmerlin/merlin-server/space/app"
	spaceappApp "github.com/openmerlin/merlin-server/spaceapp/app"
	userapp "github.com/openmerlin/merlin-server/user/app"
	"github.com/openmerlin/merlin-server/user/controller"
	userrepo "github.com/openmerlin/merlin-server/user/domain/repository"
)

type allServices struct {
	userApp  userapp.UserService
	userRepo userrepo.User

	orgApp            orgapp.OrgService
	orgCertificateApp orgapp.OrgCertificateService

	permissionApp commonapp.ResourcePermissionAppService

	sessionApp sessionapp.SessionAppService

	resourceApp coderepoapp.ResourceAppService
	codeRepoApp coderepoapp.CodeRepoAppService

	operationLog          middleware.OperationLog
	securityLog           middleware.SecurityLog
	userMiddleWare        middleware.UserMiddleWare
	rateLimiterMiddleWare middleware.RateLimiter
	privacyCheck          middleware.PrivacyCheck
	tokenMiddleWare       middleware.TokenMiddleWare

	npuGatekeeper orgapp.PrivilegeOrg
	disable       orgapp.PrivilegeOrg
	modelApp      modelapp.ModelAppService

	datasetApp datasetapp.DatasetAppService

	spaceApp spaceapp.SpaceAppService

	spaceappApp spaceappApp.SpaceappAppService

	activityApp activityapp.ActivityAppService

	modelSpace spaceapp.ModelSpaceAppService

	spaceVariable spaceapp.SpaceVariableService

	spaceSecret spaceapp.SpaceSecretService

	computilityApp computilityapp.ComputilityInternalAppService

	privacyClear controller.PrivacyClear
}

func initServices(cfg *config.Config) (services allServices, err error) {
	initUser(cfg, &services)

	// initOrg depends on initUser
	if err = initOrg(cfg, &services); err != nil {
		return
	}

	// initSession depends on initUser
	initSession(cfg, &services)

	if err = initCodeRepo(cfg, &services); err != nil {
		return
	}

	// initModel depends on initCodeRepo and initOrg
	if err = initModel(cfg, &services); err != nil {
		return
	}

	if err = initDataset(cfg, &services); err != nil {
		return
	}

	if err = initComputilityApp(cfg, &services); err != nil {
		return
	}

	if err = initSpaceApp(cfg, &services); err != nil {
		return
	}

	// initSpace depends on initCodeRepo and initOrg and initSpaceApp
	if err = initSpace(cfg, &services); err != nil {
		return
	}

	if err = initActivity(cfg, &services); err != nil {
		return
	}

	// initResource depends on initModel and initSpace
	initResource(&services)

	return
}
