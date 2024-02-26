/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/models/app"
	userapp "github.com/openmerlin/merlin-server/user/app"
)

// AddRouteForModelWebController adds a router for the ModelWebController with the given middleware.
func AddRouteForModelWebController(
	r *gin.RouterGroup,
	s app.ModelAppService,
	m middleware.UserMiddleWare,
	l middleware.OperationLog,
	u userapp.UserService,
) {
	ctl := ModelWebController{
		ModelController: ModelController{
			appService:     s,
			userMiddleWare: m,
			user:           u,
		},
	}

	addRouteForModelController(r, &ctl.ModelController, l)

	r.GET("/v1/model/:owner/:name", m.Optional, ctl.Get)
	r.GET("/v1/model/:owner", m.Optional, ctl.List)
	r.GET("/v1/model", m.Optional, ctl.ListGlobal)
}

// ModelWebController is a struct that holds the app service for model web operations.
type ModelWebController struct {
	ModelController
}

// @Summary  Get
// @Description  get model
// @Tags     ModelWeb
// @Param    owner  path  string  true  "owner of model"
// @Param    name   path  string  true  "name of model"
// @Accept   json
// @Success  200  {object}  modelDetail
// @Router   /v1/model/{owner}/{name} [get]
func (ctl *ModelWebController) Get(ctx *gin.Context) {
	index, err := ctl.parseIndex(ctx)
	if err != nil {
		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	dto, err := ctl.appService.GetByName(user, &index)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	detail := modelDetail{
		Liked:    true,
		ModelDTO: &dto,
	}

	if user != nil {
		//TODO check whether user like the model
	}

	if avatar, err := ctl.user.GetUserAvatarId(index.Owner); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		detail.AvatarId = avatar.AvatarId

		commonctl.SendRespOfGet(ctx, &detail)
	}
}

// @Summary  List
// @Description  list model
// @Tags     ModelWeb
// @Param    owner           path   string  true   "owner of model"
// @Param    name            query  string  false  "name of model"
// @Param    count           query  bool    false  "whether to calculate the total"
// @Param    sort_by         query  string  false  "sort types: most_likes, alphabetical, most_downloads, recently_updated, recently_created"
// @Param    page_num        query  int     false  "page num which starts from 1"
// @Param    count_per_page  query  int     false  "count per page"
// @Accept   json
// @Success  200  {object}  userModelsInfo
// @Router   /v1/model/{owner} [get]
func (ctl *ModelWebController) List(ctx *gin.Context) {
	var req reqToListUserModels
	if err := ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	cmd.Owner, err = primitive.NewAccount(ctx.Param("owner"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	dto, err := ctl.appService.List(user, &cmd)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	result := userModelsInfo{
		Owner:     cmd.Owner.Account(),
		ModelsDTO: &dto,
	}

	if avatar, err := ctl.user.GetUserAvatarId(cmd.Owner); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		result.AvatarId = avatar.AvatarId

		commonctl.SendRespOfGet(ctx, &result)
	}
}

// @Summary  ListGlobal
// @Description  list global public model
// @Tags     ModelWeb
// @Param    name            query  string  false  "name of model"
// @Param    task            query  string  false  "task label"
// @Param    others          query  string  false  "other labels, separate multiple each ones with commas"
// @Param    license         query  string  false  "license label"
// @Param    frameworks      query  string  false  "framework labels, separate multiple each ones with commas"
// @Param    count           query  bool    false  "whether to calculate the total"
// @Param    sort_by         query  string  false  "sort types: most_likes, alphabetical, most_downloads, recently_updated, recently_created"
// @Param    page_num        query  int     false  "page num which starts from 1"
// @Param    count_per_page  query  int     false  "count per page"
// @Accept   json
// @Success  200  {object}  modelsInfo
// @Router   /v1/model [get]
func (ctl *ModelWebController) ListGlobal(ctx *gin.Context) {
	var req reqToListGlobalModels
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

	result, err := ctl.appService.List(user, &cmd)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	if v, err := ctl.setAvatars(&result); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, v)
	}
}

func (ctl *ModelWebController) setAvatars(dto *app.ModelsDTO) (modelsInfo, error) {
	ms := dto.Models

	// get avatars

	v := map[string]bool{}
	for i := range ms {
		v[ms[i].Owner] = true
	}

	accounts := make([]primitive.Account, len(v))

	i := 0
	for k := range v {
		accounts[i] = primitive.CreateAccount(k)
		i++
	}

	avatars, err := ctl.user.GetUsersAvatarId(accounts)
	if err != nil {
		return modelsInfo{}, err
	}

	// set avatars

	am := map[string]string{}
	for i := range avatars {
		item := &avatars[i]

		am[item.Name] = item.AvatarId
	}

	infos := make([]modelInfo, len(ms))
	for i := range ms {
		item := &ms[i]

		infos[i] = modelInfo{
			AvatarId:     am[item.Owner],
			ModelSummary: item,
		}
	}

	return modelsInfo{
		Total:  dto.Total,
		Models: infos,
	}, nil
}
