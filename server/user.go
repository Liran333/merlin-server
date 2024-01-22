package server

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/common/infrastructure/gitea"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/config"
	"github.com/openmerlin/merlin-server/infrastructure/giteauser"

	"github.com/openmerlin/merlin-server/user/app"
	"github.com/openmerlin/merlin-server/user/controller"
	usergit "github.com/openmerlin/merlin-server/user/infrastructure/git"
	userrepoimpl "github.com/openmerlin/merlin-server/user/infrastructure/repositoryimpl"
)

func initUser(cfg *config.Config, services *allServices) {
	git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))

	token := userrepoimpl.NewTokenRepo(postgresql.DAO(cfg.User.Tables.Token))

	services.userRepo = userrepoimpl.NewUserRepo(postgresql.DAO(cfg.User.Tables.User))

	services.userApp = app.NewUserService(services.userRepo, git, token)
}

func setRouterOfUser(v1 *gin.RouterGroup, cfg *config.Config, services *allServices) {
	controller.AddRouterForUserController(
		v1, services.userApp, services.userRepo, services.userMiddleWare,
	)
}
