/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides functionality for managing the application's controllers.
package controller

import (
	"github.com/gin-gonic/gin"

	activityapp "github.com/openmerlin/merlin-server/activity/app"
	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/models/app"
	userapp "github.com/openmerlin/merlin-server/user/app"
)

// AddRouteForModelRestfulController adds a router for the ModelRestfulController with the given middleware.
func AddRouteForModelRestfulController(
	r *gin.RouterGroup,
	s app.ModelAppService,
	m middleware.UserMiddleWare,
	l middleware.OperationLog,
	sl middleware.SecurityLog,
	u userapp.UserService,
	rl middleware.RateLimiter,
	p middleware.PrivacyCheck,
	a activityapp.ActivityAppService,
) {
	ctl := ModelRestfulController{
		ModelController: ModelController{
			appService:     s,
			userMiddleWare: m,
			user:           u,
			activity:       a,
		},
	}

	addRouteForModelController(r, &ctl.ModelController, l, sl)

	r.GET("/v1/model/:owner/:name", p.CheckOwner, m.Optional, rl.CheckLimit, ctl.Get)
	r.GET("/v1/model", m.Optional, rl.CheckLimit, ctl.List)
}

// ModelRestfulController is a struct that holds the app service for model restful operations.
type ModelRestfulController struct {
	ModelController
}

// @Summary  Get
// @Description  get model
// @Tags     ModelRestful
// @Param    owner  path  string  true  "owner of model" MaxLength(40)
// @Param    name   path  string  true  "name of model" MaxLength(100)
// @Accept   json
// @Security Bearer
// @Success  200  {object}  commonctl.ResponseData{data=app.ModelDTO,msg=string,code=string}
// @Router   /v1/model/{owner}/{name} [get]
func (ctl *ModelRestfulController) Get(ctx *gin.Context) {
	index, err := ctl.parseIndex(ctx)
	if err != nil {
		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	dto, err := ctl.appService.GetByName(ctx.Request.Context(), user, &index)
	if err != nil {
		commonctl.SendError(ctx, err)
		return
	}

	liked := false

	modelId, _ := primitive.NewIdentity(dto.Id)
	if user != nil {
		liked, err = ctl.activity.HasLike(user, modelId)
		if err != nil {
			commonctl.SendError(ctx, err)
			return
		}
	}

	detail := modelDetail{
		Liked:    liked,
		ModelDTO: &dto,
	}

	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, &detail)
	}
}

// @Summary  List
// @Description  list global public model
// @Tags     ModelRestful
// @Param    name            query  string  false  "name of model" MaxLength(100)
// @Param    task            query  string  false  "task label" MaxLength(100)
// @Param    owner           query  string  true   "owner of model" MaxLength(40)
// @Param    others          query  string  false  "other labels, separate multiple each ones with commas" MaxLength(100)
// @Param    license         query  string  false  "license label" MaxLength(40)
// @Param    frameworks      query  string  false  "framework labels, separate multiple each ones with commas" MaxLength(100)
// @Param    count           query  bool    false  "whether to calculate the total" Enums(true, false)
// @Param    sort_by         query  string  false  "sort types: most_likes, alphabetical, most_downloads, recently_updated, recently_created" Enums(most_likes, alphabetical,most_downloads,recently_updated,recently_created)
// @Param    page_num        query  int     false  "page num which starts from 1" Mininum(1)
// @Param    count_per_page  query  int     false  "count per page" MaxCountPerPage(100)
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=app.ModelsDTO,msg=string,code=string}
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

	if result, err := ctl.appService.List(ctx.Request.Context(), user, &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, result)
	}
}
