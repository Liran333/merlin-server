package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/coderepo/app"
	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
)

func AddRouterForCodeRepoController(
	rg *gin.RouterGroup,
	cr app.CodeRepoAppService,
	m middleware.UserMiddleWare,
	rl middleware.RateLimiter) {
	ctl := CodeRepoController{
		codeRepo:       cr,
		userMiddleWare: m,
	}

	rg.GET("/v1/exists/:owner/:name", m.Read, rl.CheckLimit, ctl.Get)
}

type CodeRepoController struct {
	codeRepo       app.CodeRepoAppService
	userMiddleWare middleware.UserMiddleWare
}

// @Summary  Check
// @Description  check whether the repo exists
// @Tags     CodeRepo
// @Param    owner  path  string  true  "owner of repo"
// @Param    name   path  string  true  "name of repo"
// @Accept   json
// @Success  200  {object} bool
// @Router   /v1/exists/:owner/:name [get]
func (ctl *CodeRepoController) Get(ctx *gin.Context) {
	codeRepo, err := ToCmdToCheckRepoExists(ctx)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
	}

	data, err := ctl.codeRepo.IsRepoExist(codeRepo)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, data)
	}
}
