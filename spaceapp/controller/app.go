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
	"github.com/openmerlin/merlin-server/space/domain"
	"github.com/openmerlin/merlin-server/spaceapp/app"
)

// AddRouterForSpaceappWebController adds a router for the SpaceAppWebController to the given gin.RouterGroup.
func AddRouterForSpaceappWebController(
	r *gin.RouterGroup,
	s app.SpaceappAppService,
	m middleware.UserMiddleWare,
	l middleware.RateLimiter,
) {
	ctl := SpaceAppWebController{
		appService:          s,
		userMiddleWare:      m,
		rateLimitMiddleWare: l,
	}

	r.GET("/v1/space-app/:owner/:name", m.Optional, l.CheckLimit, ctl.Get)
}

// SpaceAppWebController is a struct that represents the web controller for the space app.
type SpaceAppWebController struct {
	appService          app.SpaceappAppService
	userMiddleWare      middleware.UserMiddleWare
	rateLimitMiddleWare middleware.RateLimiter
}

// @Summary  Get
// @Description  get space app
// @Tags     SpaceAppWeb
// @Param    owner  path  string  true  "owner of space"
// @Param    name   path  string  true  "name of space"
// @Accept   json
// @Success  200  {object}  app.SpaceAppDTO
// @Router   /v1/space-app/{owner}/{name} [get]
func (ctl *SpaceAppWebController) Get(ctx *gin.Context) {
	index, err := ctl.parseIndex(ctx)
	if err != nil {
		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	if dto, err := ctl.appService.GetByName(user, &index); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, &dto)
	}
}

func (ctl *SpaceAppWebController) parseIndex(ctx *gin.Context) (index domain.SpaceIndex, err error) {
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
