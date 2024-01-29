package server

import (
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/config"

	userapp "github.com/openmerlin/merlin-server/user/app"
	userrepo "github.com/openmerlin/merlin-server/user/domain/repository"

	orgapp "github.com/openmerlin/merlin-server/organization/app"

	modelapp "github.com/openmerlin/merlin-server/models/app"

	spaceapp "github.com/openmerlin/merlin-server/space/app"

	sessionapp "github.com/openmerlin/merlin-server/session/app"

	coderepoapp "github.com/openmerlin/merlin-server/coderepo/app"
)

type allServices struct {
	userApp  userapp.UserService
	userRepo userrepo.User

	orgApp     orgapp.OrgService
	permission orgapp.Permission

	sessionApp sessionapp.SessionAppService

	codeRepoApp     coderepoapp.CodeRepoAppService
	codeRepoFileApp coderepoapp.CodeRepoFileAppService

	userMiddleWare middleware.UserMiddleWare

	modelApp modelapp.ModelAppService

	spaceApp spaceapp.SpaceAppService
}

func initServices(cfg *config.Config) (services allServices, err error) {
	if err = initCodeRepo(cfg, &services); err != nil {
		return
	}

	if err = initModel(cfg, &services); err != nil {
		return
	}

	if err = initSpace(cfg, &services); err != nil {
		return
	}

	initUser(cfg, &services)

	// initOrg depends on initUser
	initOrg(cfg, &services)

	// initSession depends on initUser
	initSession(cfg, &services)

	return
}
