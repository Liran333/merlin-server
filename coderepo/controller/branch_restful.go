package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/coderepo/app"
	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
)

func AddRouteForBranchRestfulController(
	r *gin.RouterGroup,
	s app.BranchAppService,
	m middleware.UserMiddleWare,
	l middleware.OperationLog,
) {
	ctl := BranchRestfulController{
		userMiddleWare: m,
		appService:     s,
	}

	r.POST("/v1/branch/:type/:owner/:repo", m.Optional, l.Write, ctl.Create)
	r.DELETE("/v1/branch/:type/:owner/:repo/:branch", m.Optional, l.Write, ctl.Delete)
}

type BranchRestfulController struct {
	userMiddleWare middleware.UserMiddleWare
	appService     app.BranchAppService
}

// @Summary  CreateBranch
// @Description  create repo branch
// @Tags     BranchRestful
// @Param    type  path  string  true  "type of space/model"
// @Param    owner  path  string  true  "owner of space/model"
// @Param    repo  path  string  true  "name of space/model"
// @Param    body     body restfulReqToCreateBranch true  "restfulReqToCreateBranch"
// @Accept   json
// @Success  201   {object}  app.BranchCreateDTO
// @Router   /v1/branch/{type}/{owner}/{repo} [post]
func (ctl *BranchRestfulController) Create(ctx *gin.Context) {
	var req restfulReqToCreateBranch
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	middleware.SetAction(ctx, req.action(ctx))

	cmd, err := req.toCmd(ctx)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	v, err := ctl.appService.Create(user, &cmd)
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, &v)
	}
}

// @Summary  DeleteBranch
// @Description  delete repo branch
// @Tags     BranchRestful
// @Param    type  path  string  true  "repo type"
// @Param    owner  path  string  true  "repo owner"
// @Param    repo  path  string  true  "repo name"
// @Param    branch  path  string  true  "branch name"
// @Accept   json
// @Success  204
// @Router   /v1/branch/{type}/{owner}/{repo}/{branch} [delete]
func (ctl *BranchRestfulController) Delete(ctx *gin.Context) {
	cmd, err := toBanchDeleteCmd(ctx)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	middleware.SetAction(ctx, fmt.Sprintf("delete branch %s/%s/%s",
		ctx.Param("owner"), ctx.Param("repo"), ctx.Param("branch")))

	user := ctl.userMiddleWare.GetUser(ctx)

	err = ctl.appService.Delete(user, &cmd)
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfDelete(ctx)
	}
}
