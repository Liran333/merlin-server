package server

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/coderepo/app"
	"github.com/openmerlin/merlin-server/coderepo/controller"
	"github.com/openmerlin/merlin-server/coderepo/infrastructure/coderepoadapter"
	"github.com/openmerlin/merlin-server/coderepo/infrastructure/coderepofileadapter"
	"github.com/openmerlin/merlin-server/common/infrastructure/gitea"
)

func initCodeRepo(services *allServices) {
	services.codeRepoApp = app.NewCodeRepoAppService(
		coderepoadapter.NewRepoAdapter(gitea.Client()),
	)

	services.codeRepoFileApp = app.NewCodeRepoFileAppService(
		coderepofileadapter.NewCodeRepoFileAdapter(gitea.Client()),
	)
}

func setRouteOfCodeRepoFile(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouterForCodeRepoFileController(
		rg,
		services.codeRepoFileApp,
		services.userMiddleWare,
	)
}
