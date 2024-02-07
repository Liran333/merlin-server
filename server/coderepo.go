package server

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/coderepo/app"
	"github.com/openmerlin/merlin-server/coderepo/controller"
	coderepoprimtive "github.com/openmerlin/merlin-server/coderepo/domain/primitive"
	"github.com/openmerlin/merlin-server/coderepo/infrastructure/branchclientadapter"
	"github.com/openmerlin/merlin-server/coderepo/infrastructure/branchrepositoryadapter"
	"github.com/openmerlin/merlin-server/coderepo/infrastructure/coderepoadapter"
	"github.com/openmerlin/merlin-server/coderepo/infrastructure/coderepofileadapter"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/common/infrastructure/gitea"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/config"

	modelapp "github.com/openmerlin/merlin-server/models/app"
	modeldomain "github.com/openmerlin/merlin-server/models/domain"

	spaceapp "github.com/openmerlin/merlin-server/space/app"
	spaceomain "github.com/openmerlin/merlin-server/space/domain"
)

func initCodeRepo(cfg *config.Config, services *allServices) error {
	err := branchrepositoryadapter.Init(postgresql.DB(), &cfg.CodeRepo.Tables)
	if err != nil {
		return err
	}

	services.codeRepoApp = app.NewCodeRepoAppService(
		coderepoadapter.NewRepoAdapter(gitea.Client()),
	)

	services.codeRepoFileApp = app.NewCodeRepoFileAppService(
		coderepofileadapter.NewCodeRepoFileAdapter(gitea.Client()),
	)

	return nil
}

func setRouterOfCodeRepoFile(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouterForCodeRepoFileController(
		rg,
		services.codeRepoFileApp,
		services.userMiddleWare,
	)
}

func setRouterOfBranchRestful(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouteForBranchRestfulController(
		rg,
		app.NewBranchAppService(
			services.permission,
			branchrepositoryadapter.BranchAdapter(),
			branchclientadapter.NewBranchClientAdapter(gitea.Client()),
			NewCheckRepoAdapter(services),
		),
		services.userMiddleWare,
		services.operationLog,
	)
}

func NewCheckRepoAdapter(services *allServices) *checkRepoAdapter {
	return &checkRepoAdapter{
		modelApp: services.modelApp,
		spaceApp: services.spaceApp,
	}
}

type checkRepoAdapter struct {
	modelApp modelapp.ModelAppService
	spaceApp spaceapp.SpaceAppService
}

func (a *checkRepoAdapter) CheckRepo(t coderepoprimtive.RepoType, account primitive.Account, name primitive.MSDName) (err error) {
	if t.IsModel() {
		_, err = a.modelApp.GetByName(account, &modeldomain.ModelIndex{Owner: account, Name: name})

	} else if t.IsSpace() {
		_, err = a.spaceApp.GetByName(account, &spaceomain.SpaceIndex{Owner: account, Name: name})

	}

	return
}
