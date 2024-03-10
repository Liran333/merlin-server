/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package server

import (
	"github.com/gin-gonic/gin"
	redisdb "github.com/opensourceways/redis-lib"

	"github.com/openmerlin/merlin-server/common/domain/crypto"
	"github.com/openmerlin/merlin-server/common/infrastructure/gitea"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/config"
	"github.com/openmerlin/merlin-server/infrastructure/giteauser"
	sessionapp "github.com/openmerlin/merlin-server/session/app"
	"github.com/openmerlin/merlin-server/session/infrastructure/loginrepositoryadapter"
	"github.com/openmerlin/merlin-server/session/infrastructure/oidcimpl"
	"github.com/openmerlin/merlin-server/session/infrastructure/sessionrepositoryadapter"
	"github.com/openmerlin/merlin-server/user/app"
	"github.com/openmerlin/merlin-server/user/controller"
	usergit "github.com/openmerlin/merlin-server/user/infrastructure/git"
	userrepoimpl "github.com/openmerlin/merlin-server/user/infrastructure/repositoryimpl"
)

func initUser(cfg *config.Config, services *allServices) {
	git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))

	token := userrepoimpl.NewTokenRepo(postgresql.DAO(cfg.User.Tables.Token))

	services.userRepo = userrepoimpl.NewUserRepo(postgresql.DAO(cfg.User.Tables.User),
		crypto.NewEncryption(cfg.User.Key))

	session := sessionapp.NewSessionClearAppService(
		loginrepositoryadapter.LoginAdapter(),
		sessionrepositoryadapter.NewSessionAdapter(redisdb.DAO()),
	)

	services.userApp = app.NewUserService(
		services.userRepo,
		git,
		token,
		loginrepositoryadapter.LoginAdapter(),
		oidcimpl.NewAuthingUser(),
		session,
		cfg.User,
	)
}

func setRouterOfUser(v1 *gin.RouterGroup, cfg *config.Config, services *allServices) {
	controller.AddRouterForUserController(
		v1,
		services.userApp,
		services.userRepo,
		services.operationLog,
		services.userMiddleWare,
		services.rateLimiterMiddleWare,
		services.privacyCheck,
	)
}

func setInternalRouterOfUser(v1 *gin.RouterGroup, cfg *config.Config, services *allServices) {
	controller.AddRouterForUserInternalController(
		v1,
		services.userApp,
		services.userMiddleWare,
	)
}
