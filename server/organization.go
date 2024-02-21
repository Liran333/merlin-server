package server

import (
	"github.com/gin-gonic/gin"

	commonapp "github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/crypto"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/config"
	"github.com/openmerlin/merlin-server/organization/app"
	"github.com/openmerlin/merlin-server/organization/controller"
	orgrepoimpl "github.com/openmerlin/merlin-server/organization/infrastructure/repositoryimpl"
	userrepoimpl "github.com/openmerlin/merlin-server/user/infrastructure/repositoryimpl"
)

// initOrg depends on initUser
func initOrg(cfg *config.Config, services *allServices) {
	org := userrepoimpl.NewUserRepo(postgresql.DAO(cfg.User.Tables.User), crypto.NewEncryption(cfg.User.Key))

	orgMember := orgrepoimpl.NewMemberRepo(postgresql.DAO(cfg.Org.Tables.Member))

	invitation := orgrepoimpl.NewInviteRepo(postgresql.DAO(cfg.Org.Tables.Invite))

	permission := app.NewPermService(&cfg.Permission, orgMember)

	services.permissionApp = commonapp.NewResourcePermissionAppService(permission)

	services.orgApp = app.NewOrgService(
		services.userApp, org, orgMember,
		invitation, permission, &cfg.Org,
	)
}

func setRouterOfOrg(v1 *gin.RouterGroup, cfg *config.Config, services *allServices) {
	controller.AddRouterForOrgController(
		v1,
		services.orgApp,
		services.userApp,
		services.operationLog,
		services.userMiddleWare,
	)
}
