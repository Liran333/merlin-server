/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/models/app"
	userapp "github.com/openmerlin/merlin-server/user/app"
)

// AddRouteForModelRestfulController adds a router for the ModelRestfulController with the given middleware.
func AddRouteForModelRestfulController(
	r *gin.RouterGroup,
	s app.ModelAppService,
	m middleware.UserMiddleWare,
	l middleware.OperationLog,
	u userapp.UserService,
) {
	ctl := ModelRestfulController{
		ModelController: ModelController{
			appService:     s,
			userMiddleWare: m,
			user:           u,
		},
	}

	addRouteForModelController(r, &ctl.ModelController, l)

	r.GET("/v1/model/:owner/:name", m.Optional, ctl.Get)
	r.GET("/v1/model", m.Optional, ctl.List)
}

// ModelRestfulController is a struct that holds the app service for model restful operations.
type ModelRestfulController struct {
	ModelController
}

// @Summary  Get
// @Description  get model
// @Tags     ModelRestful
// @Param    owner  path  string  true  "owner of model"
// @Param    name   path  string  true  "name of model"
// @Accept   json
// @Success  200  {object}  app.ModelDTO
// @Router   /v1/model/{owner}/{name} [get]
func (ctl *ModelRestfulController) Get(ctx *gin.Context) {
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
// @Description  list global public model
// @Tags     ModelRestful
// @Param    name            query  string  false  "name of model"
// @Param    task            query  string  false  "task label"
// @Param    owner           query  string  true   "owner of model"
// @Param    others          query  string  false  "other labels, separate multiple each ones with commas"
// @Param    license         query  string  false  "license label"
// @Param    frameworks      query  string  false  "framework labels, separate multiple each ones with commas"
// @Param    count           query  bool    false  "whether to calculate the total"
// @Param    sort_by         query  string  false  "sort types: most_likes, alphabetical, most_downloads, recently_updated, recently_created"
// @Param    page_num        query  int     false  "page num which starts from 1"
// @Param    count_per_page  query  int     false  "count per page"
// @Accept   json
// @Success  200  {object}  app.ModelsDTO
// @Router   /v1/model [get]
func (ctl *ModelRestfulController) List(ctx *gin.Context) {
	var req restfulReqToListModels

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
