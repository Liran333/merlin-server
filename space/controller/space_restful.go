/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/space/app"
	userapp "github.com/openmerlin/merlin-server/user/app"
)

// AddRouteForSpaceRestfulController adds routes to the given router group for the SpaceRestfulController.
func AddRouteForSpaceRestfulController(
	r *gin.RouterGroup,
	s app.SpaceAppService,
	m middleware.UserMiddleWare,
	l middleware.OperationLog,
	rl middleware.RateLimiter,
	u userapp.UserService,
) {
	ctl := SpaceRestfulController{
		SpaceController: SpaceController{
			appService:     	 s,
			userMiddleWare: 	 m,
			rateLimitMiddleWare: rl,
			user:           	 u,
		},
	}

	addRouteForSpaceController(r, &ctl.SpaceController, l, rl)

	r.GET("/v1/space/:owner/:name", m.Optional, ctl.Get)
	r.GET("/v1/space", m.Optional, ctl.List)
}

// SpaceRestfulController is a struct that holds the necessary dependencies for handling space-related operations.
type SpaceRestfulController struct {
	SpaceController
}

// @Summary  Get
// @Description  get space
// @Tags     SpaceRestful
// @Param    owner  path  string  true  "owner of space"
// @Param    name   path  string  true  "name of space"
// @Accept   json
// @Success  200  {object}  app.SpaceDTO
// @Router   /v1/space/{owner}/{name} [get]
func (ctl *SpaceRestfulController) Get(ctx *gin.Context) {
	index, err := ctl.parseIndex(ctx)
	if err != nil {
		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	dto, err := ctl.appService.GetByName(user, &index)
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, &dto)
	}
}

// @Summary  List
// @Description  list global public space
// @Tags     SpaceRestful
// @Param    name            query  string  false  "name of space"
// @Param    task            query  string  false  "task label"
// @Param    owner           query  string  true   "owner of space"
// @Param    others          query  string  false  "other labels, separate multiple each ones with commas"
// @Param    license         query  string  false  "license label"
// @Param    frameworks      query  string  false  "framework labels, separate multiple each ones with commas"
// @Param    count           query  bool    false  "whether to calculate the total"
// @Param    sort_by         query  string  false  "sort types: most_likes, alphabetical, most_downloads, recently_updated, recently_created"
// @Param    page_num        query  int     false  "page num which starts from 1"
// @Param    count_per_page  query  int     false  "count per page"
// @Accept   json
// @Success  200  {object}  app.SpacesDTO
// @Router   /v1/space [get]
func (ctl *SpaceRestfulController) List(ctx *gin.Context) {
	var req restfulReqToListSpaces

	if err := ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	if result, err := ctl.appService.List(user, &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, result)
	}
}
