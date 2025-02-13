/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package server provides functionality for setting up and configuring a server for handling code repo operations.
package server

import (
	"github.com/gin-gonic/gin"
	redisdb "github.com/opensourceways/redis-lib"

	"github.com/openmerlin/merlin-server/common/domain/crypto"
	"github.com/openmerlin/merlin-server/common/infrastructure/email"
	"github.com/openmerlin/merlin-server/common/infrastructure/gitea"
	"github.com/openmerlin/merlin-server/common/infrastructure/obs"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/config"
	"github.com/openmerlin/merlin-server/infrastructure/giteauser"
	orgrepoimpl "github.com/openmerlin/merlin-server/organization/infrastructure/repositoryimpl"
	sessionapp "github.com/openmerlin/merlin-server/session/app"
	"github.com/openmerlin/merlin-server/session/infrastructure/loginrepositoryadapter"
	"github.com/openmerlin/merlin-server/session/infrastructure/oidcimpl"
	"github.com/openmerlin/merlin-server/session/infrastructure/sessionrepositoryadapter"
	"github.com/openmerlin/merlin-server/space/infrastructure/obsadapter"
	"github.com/openmerlin/merlin-server/user/app"
	"github.com/openmerlin/merlin-server/user/controller"
	"github.com/openmerlin/merlin-server/user/infrastructure/emailadapter"
	usergit "github.com/openmerlin/merlin-server/user/infrastructure/git"
	userrepoimpl "github.com/openmerlin/merlin-server/user/infrastructure/repositoryimpl"
)

func initUser(cfg *config.Config, services *allServices) {
	git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))

	token := userrepoimpl.NewTokenRepo(postgresql.DAO(cfg.User.Domain.Tables.Token))

	services.userRepo = userrepoimpl.NewUserRepo(postgresql.DAO(cfg.User.Domain.Tables.User),
		crypto.NewEncryption(cfg.User.Domain.Key))

	member := orgrepoimpl.NewMemberRepo(postgresql.DAO(cfg.Org.Domain.Tables.Member))

	session := sessionapp.NewSessionClearAppService(
		loginrepositoryadapter.LoginAdapter(),
		sessionrepositoryadapter.NewSessionAdapter(redisdb.DAO()),
	)

	services.userApp = app.NewUserService(
		services.userRepo,
		member,
		git,
		token,
		loginrepositoryadapter.LoginAdapter(),
		oidcimpl.NewAuthingUser(),
		session,
		&cfg.User.Domain,
		obsadapter.NewClient(obs.Client()),
		emailadapter.NewEmailImpl(email.GetEmailInst(), cfg.User.Domain.AuditEmail, cfg.User.Domain.MailTemplates),
	)
}

func setRouterOfUser(v1 *gin.RouterGroup, cfg *config.Config, services *allServices) {
	controller.AddRouterForUserController(
		v1,
		services.userApp,
		services.userRepo,
		services.operationLog,
		services.securityLog,
		services.userMiddleWare,
		services.rateLimiterMiddleWare,
		services.privacyCheck,
		services.disable,
		services.privacyClear,
	)
}

func setInternalRouterOfUser(v1 *gin.RouterGroup, cfg *config.Config, services *allServices) {
	controller.AddRouterForUserInternalController(
		v1,
		services.userApp,
		services.userMiddleWare,
	)
}
