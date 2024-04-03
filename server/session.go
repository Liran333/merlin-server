/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package server provides functionality for setting up and configuring a server for handling code repo operations.
package server

import (
	"github.com/gin-gonic/gin"
	redisdb "github.com/opensourceways/redis-lib"

	"github.com/openmerlin/merlin-server/config"
	sessionapp "github.com/openmerlin/merlin-server/session/app"
	"github.com/openmerlin/merlin-server/session/controller"
	sessionctl "github.com/openmerlin/merlin-server/session/controller"
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

func setRouterOfSession(rg *gin.RouterGroup, services *allServices, cfg *sessionctl.Config) {
	controller.AddRouterForSessionController(
		rg, services.sessionApp, services.operationLog, services.userMiddleWare, cfg,
	)
}

func setRouterOfSessionInternal(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouterForSessionInternalController(
		rg, services.sessionApp, services.operationLog, services.userMiddleWare,
	)
}
