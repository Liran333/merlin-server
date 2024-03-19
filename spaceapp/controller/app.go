/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides the controllers for handling HTTP requests and managing the application's business logic.
package controller

import (
	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	spacedomain "github.com/openmerlin/merlin-server/space/domain"
	"github.com/openmerlin/merlin-server/spaceapp/app"
)

func addRouterForSpaceappController(
	r *gin.RouterGroup,
	ctl *SpaceAppController,
	m middleware.UserMiddleWare,
	l middleware.RateLimiter,
) {

	r.POST("/v1/space-app/:owner/:name/restart", m.Optional, l.CheckLimit, ctl.Restart)

}

// SpaceAppController is a struct that represents the  controller for the space app.
type SpaceAppController struct {
	appService          app.SpaceappAppService
	userMiddleWare      middleware.UserMiddleWare
	tokenMiddleWare     middleware.TokenMiddleWare
	rateLimitMiddleWare middleware.RateLimiter
}

func (ctl *SpaceAppController) parseIndex(ctx *gin.Context) (index spacedomain.SpaceIndex, err error) {
	index.Owner, err = primitive.NewAccount(ctx.Param("owner"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	index.Name, err = primitive.NewMSDName(ctx.Param("name"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
	}

	return
}

// @Summary  Post
// @Description  restart space app
// @Tags     SpaceApp
// @Param    owner  path  string  true  "owner of space"
// @Param    name   path  string  true  "name of space"
// @Accept   json
// @Success  201   {object}  commonctl.ResponseData
// @Router   /v1/space-app/{owner}/{name}/restart [post]
func (ctl *SpaceAppController) Restart(ctx *gin.Context) {
	index, err := ctl.parseIndex(ctx)
	if err != nil {
		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	if err := ctl.appService.RestartSpaceApp(user, &index); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, "successfully")
	}
}
