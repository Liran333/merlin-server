package server

import (
	coderepoapp "github.com/openmerlin/merlin-server/coderepo/app"
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

	orgApp     orgapp.OrgService
	permission orgapp.Permission

	sessionApp sessionapp.SessionAppService

	codeRepoApp     coderepoapp.CodeRepoAppService
	codeRepoFileApp coderepoapp.CodeRepoFileAppService

	operationLog   middleware.OperationLog
	userMiddleWare middleware.UserMiddleWare

	modelApp modelapp.ModelAppService

	spaceApp spaceapp.SpaceAppService
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

	// initSpace depends on initCodeRepo and initOrg
	if err = initSpace(cfg, &services); err != nil {
		return
	}

	if err = initSpaceApp(cfg, &services); err != nil {
		return
	}

	return
}
