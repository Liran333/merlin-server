/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides the controllers for handling restful requests and converting them into commands
package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/coderepo/app"
	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
)

// AddRouteForBranchRestfulController adds routes for BranchRestfulController to the given router group.
func AddRouteForBranchRestfulController(
	r *gin.RouterGroup,
	s app.BranchAppService,
	m middleware.UserMiddleWare,
	l middleware.OperationLog,
	rl middleware.RateLimiter,
) {
	ctl := BranchRestfulController{
		userMiddleWare: m,
		appService:     s,
	}

	r.POST("/v1/branch/:type/:owner/:repo", m.Write, l.Write, rl.CheckLimit, ctl.Create)
	r.DELETE("/v1/branch/:type/:owner/:repo/:branch", m.Write, l.Write, rl.CheckLimit, ctl.Delete)
}

// BranchRestfulController is a struct that holds user middleware and app service for branch operations.
type BranchRestfulController struct {
	userMiddleWare middleware.UserMiddleWare
	appService     app.BranchAppService
}

// @Summary  CreateBranch
// @Description  create repo branch
// @Tags     BranchRestful
// @Param    type  path  string  true  "type of space/model" Enums(space, model)
// @Param    owner  path  string  true  "owner of space/model" MaxLength(40)
// @Param    repo  path  string  true  "name of space/model" MaxLength(100)
// @Param    body     body restfulReqToCreateBranch true  "restfulReqToCreateBranch"
// @Accept   json
// @Security Bearer
// @Success  201   {object}  commonctl.ResponseData{data=app.BranchCreateDTO,msg=string,code=string}
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
// @Param    type  path  string  true  "repo type" Enums(space, model)
// @Param    owner  path  string  true  "repo owner" MaxLength(40)
// @Param    repo  path  string  true  "repo name" MaxLength(100)
// @Param    branch  path  string  true  "branch name" MaxLength(100)
// @Accept   json
// @Security Bearer
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
