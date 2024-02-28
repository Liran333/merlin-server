/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package server

import (
	"github.com/gin-gonic/gin"
	redisdb "github.com/opensourceways/redis-lib"

	"github.com/openmerlin/merlin-server/config"
	sessionapp "github.com/openmerlin/merlin-server/session/app"
	"github.com/openmerlin/merlin-server/session/controller"
	"github.com/openmerlin/merlin-server/session/infrastructure/csrftokenrepositoryadapter"
	"github.com/openmerlin/merlin-server/session/infrastructure/loginrepositoryadapter"
	"github.com/openmerlin/merlin-server/session/infrastructure/oidcimpl"
	"github.com/openmerlin/merlin-server/session/infrastructure/sessionrepositoryadapter"
)

// initSession depends on initUser
func initSession(cfg *config.Config, services *allServices) {
	services.sessionApp = sessionapp.NewSessionAppService(
		oidcimpl.NewAuthingUser(),
		services.userApp,
		cfg.Session.Domain.MaxSessionNum,
		loginrepositoryadapter.LoginAdapter(),
		csrftokenrepositoryadapter.NewCSRFTokenAdapter(redisdb.DAO()),
		sessionrepositoryadapter.NewSessionAdapter(redisdb.DAO()),
	)
}

func setRouterOfSession(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouterForSessionController(
		rg, services.sessionApp, services.operationLog, services.userMiddleWare,
	)
}

func setRouterOfSessionInternal(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouterForSessionInternalController(
		rg, services.sessionApp, services.operationLog, services.userMiddleWare,
	)
}
