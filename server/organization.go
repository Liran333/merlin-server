package server

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/config"

	"github.com/openmerlin/merlin-server/organization/app"
	"github.com/openmerlin/merlin-server/organization/controller"
	orgrepoimpl "github.com/openmerlin/merlin-server/organization/infrastructure/repositoryimpl"

	userrepoimpl "github.com/openmerlin/merlin-server/user/infrastructure/repositoryimpl"
)

// initOrg depends on initUser
func initOrg(cfg *config.Config, services *allServices) {
	org := userrepoimpl.NewUserRepo(postgresql.DAO(cfg.User.Tables.User))

	orgMember := orgrepoimpl.NewMemberRepo(postgresql.DAO(cfg.Org.Tables.Member))

	invitation := orgrepoimpl.NewInviteRepo(postgresql.DAO(cfg.Org.Tables.Invite))

	services.permission = app.NewPermService(&cfg.Permission, orgMember)

	services.orgApp = app.NewOrgService(
		services.userApp, org, orgMember,
		invitation, services.permission, &cfg.Org,
	)
}

func setRouterOfOrg(v1 *gin.RouterGroup, cfg *config.Config, services *allServices) {
	controller.AddRouterForOrgController(
		v1, services.orgApp, services.userApp, services.userMiddleWare,
	)
}
