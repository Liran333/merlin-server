/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package server

import (
	activityapp "github.com/openmerlin/merlin-server/activity/app"
	coderepoapp "github.com/openmerlin/merlin-server/coderepo/app"
	commonapp "github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/config"
	modelapp "github.com/openmerlin/merlin-server/models/app"
	orgapp "github.com/openmerlin/merlin-server/organization/app"
	sessionapp "github.com/openmerlin/merlin-server/session/app"
	spaceapp "github.com/openmerlin/merlin-server/space/app"
	userapp "github.com/openmerlin/merlin-server/user/app"
	userrepo "github.com/openmerlin/merlin-server/user/domain/repository"
)

type allServices struct {
	userApp  userapp.UserService
	userRepo userrepo.User

	orgApp orgapp.OrgService

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

	modelApp modelapp.ModelAppService

	spaceApp spaceapp.SpaceAppService

	activityApp activityapp.ActivityAppService

	modelSpace spaceapp.ModelSpaceAppService

	spaceVariable spaceapp.SpaceVariableService

	spaceSecret spaceapp.SpaceSecretService
}

func initServices(cfg *config.Config) (services allServices, err error) {
	initUser(cfg, &services)

	// initOrg depends on initUser
	initOrg(cfg, &services)

	// initSession depends on initUser
	initSession(cfg, &services)

	if err = initCodeRepo(cfg, &services); err != nil {
		return
	}

	// initModel depends on initCodeRepo and initOrg
	if err = initModel(cfg, &services); err != nil {
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
