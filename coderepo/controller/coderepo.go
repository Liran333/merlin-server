/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides the controllers for handling restful requests and converting them into commands
package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/coderepo/app"
	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
)

// AddRouterForCodeRepoController adds routes for the CodeRepoController to the specified gin.RouterGroup.
func AddRouterForCodeRepoController(
	rg *gin.RouterGroup,
	r app.ResourceAppService,
	m middleware.UserMiddleWare,
	rl middleware.RateLimiter,
) {
	ctl := CodeRepoController{
		coderepo:       r,
		userMiddleWare: m,
	}

	rg.GET("/v1/exists/:owner/:name", m.Read, rl.CheckLimit, ctl.Get)
}

// CodeRepoController is a controller that handles code repository-related operations.
type CodeRepoController struct {
	coderepo       app.ResourceAppService
	userMiddleWare middleware.UserMiddleWare
}

// @Summary  Get coderepo and check it
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

	data, err := ctl.coderepo.IsRepoExist(codeRepo)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, data)
	}
}
