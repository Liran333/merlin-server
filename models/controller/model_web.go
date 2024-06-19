/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides functionality for managing the application's controllers.
package controller

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	activityapp "github.com/openmerlin/merlin-server/activity/app"
	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/models/app"
	spaceapp "github.com/openmerlin/merlin-server/space/app"
	userapp "github.com/openmerlin/merlin-server/user/app"
)

// AddRouteForModelWebController adds a router for the ModelWebController with the given middleware.
func AddRouteForModelWebController(
	r *gin.RouterGroup,
	s app.ModelAppService,
	ms spaceapp.ModelSpaceAppService,
	m middleware.UserMiddleWare,
	l middleware.OperationLog,
	sl middleware.SecurityLog,
	u userapp.UserService,
	rl middleware.RateLimiter,
	p middleware.PrivacyCheck,
	a activityapp.ActivityAppService,
) {
	ctl := ModelWebController{
		ModelController: ModelController{
			appService:     s,
			userMiddleWare: m,
			user:           u,
			activity:       a,
		},
		modelSpaceService: ms,
	}

	addRouteForModelController(r, &ctl.ModelController, l, sl)

	r.GET("/v1/model/:owner/:name", p.CheckOwner, m.Optional, ctl.Get)
	r.GET("/v1/model/:owner", p.CheckOwner, m.Optional, ctl.List)
	r.GET("/v1/model", m.Optional, ctl.ListGlobal)
	r.GET("/v1/model/relation/:id/space", m.Optional, rl.CheckLimit, ctl.GetSpacesByModelId)

	r.PUT("/v1/model/:id/disable", ctl.ModelController.userMiddleWare.Write, l.Write, ctl.disable)
}

// ModelWebController is a struct that holds the app service for model web operations.
type ModelWebController struct {
	ModelController
	modelSpaceService spaceapp.ModelSpaceAppService
}

// @Summary  Get
// @Description  get model
// @Tags     ModelWeb
// @Param    owner  path  string  true  "owner of model" MaxLength(40)
// @Param    name   path  string  true  "name of model" MaxLength(100)
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=modelDetail,msg=string,code=string}
// @Router   /v1/model/{owner}/{name} [get]
func (ctl *ModelWebController) Get(ctx *gin.Context) {
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

	if userInfo, err := ctl.user.GetOrgOrUser(ctx.Request.Context(), user, index.Owner); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		detail.AvatarId = userInfo.AvatarId
		detail.OwnerType = userInfo.Type

		commonctl.SendRespOfGet(ctx, &detail)
	}
}

// @Summary  List
// @Description  list model
// @Tags     ModelWeb
// @Param    owner           path   string  true   "owner of model" MaxLength(40)
// @Param    name            query  string  false  "name of model" MaxLength(100)
// @Param    count           query  bool    false  "whether to calculate the total" Enums(true, false)
// @Param    sort_by         query  string  false  "sort types: most_likes, alphabetical, most_downloads, recently_updated, recently_created" Enums(most_likes, alphabetical,most_downloads,recently_updated,recently_created)
// @Param    page_num        query  int     false  "page num which starts from 1" Mininum(1)
// @Param    count_per_page  query  int     false  "count per page" MaxCountPerPage(100)
// @Accept   json
// @Success  200  {object}  userModelsInfo
// @Router   /v1/model/{owner} [get]
func (ctl *ModelWebController) List(ctx *gin.Context) {
	owner := ctx.Param("owner")

	switch owner {
	case "recommend":
		ctl.ListRecommends(ctx)
		return
	default:
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

		dto, err := ctl.appService.List(ctx.Request.Context(), user, &cmd)
		if err != nil {
			commonctl.SendError(ctx, err)

			return
		}

		result := userModelsInfo{
			Owner:     cmd.Owner.Account(),
			ModelsDTO: &dto,
		}

		if userInfo, err := ctl.user.GetOrgOrUser(ctx.Request.Context(), user, cmd.Owner); err != nil {
			commonctl.SendError(ctx, err)
		} else {
			result.AvatarId = userInfo.AvatarId
			result.OwnerType = userInfo.Type

			commonctl.SendRespOfGet(ctx, &result)
		}
	}
}

// @Summary  ListRecommends
// @Description  list recommend models
// @Tags     ModelWeb
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=modelsRecommendInfo,msg=string,code=string}
// @Router   /v1/model/recommend [get]
func (ctl *ModelWebController) ListRecommends(ctx *gin.Context) {
	user := ctl.userMiddleWare.GetUser(ctx)

	modelsDTO := ctl.appService.Recommend(ctx.Request.Context(), user)

	mrs := make([]modelRecommendInfo, 0, len(modelsDTO))

	for _, v := range modelsDTO {
		m := v
		mr := modelRecommendInfo{
			ModelDTO: &m,
		}

		mrs = append(mrs, mr)
	}

	result := modelsRecommendInfo{
		Models: mrs,
	}

	commonctl.SendRespOfGet(ctx, result)
}

// @Summary  ListGlobal
// @Description  list global public model
// @Tags     ModelWeb
// @Param    name            query  string  false  "name of model" MaxLength(100)
// @Param    task            query  string  false  "task label" MaxLength(100)
// @Param    others          query  string  false  "other labels, separate multiple each ones with commas" MaxLength(100)
// @Param    license         query  string  false  "license label" MaxLength(40)
// @Param    frameworks      query  string  false  "framework labels, separate multiple each ones with commas" MaxLength(100)
// @Param    count           query  bool    false  "whether to calculate the total" Enums(true, false)
// @Param    sort_by         query  string  false  "sort types: most_likes, alphabetical, most_downloads, recently_updated, recently_created" Enums(most_likes, alphabetical,most_downloads,recently_updated,recently_created)
// @Param    page_num        query  int     false  "page num which starts from 1" Mininum(1)
// @Param    count_per_page  query  int     false  "count per page" MaxCountPerPage(100)
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=modelsInfo,msg=string,code=string}
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

	result, err := ctl.appService.List(ctx, user, &cmd)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	if v, err := ctl.setUserInfo(ctx, &result); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, v)
	}
}

func (ctl *ModelWebController) setUserInfo(ctx context.Context, dto *app.ModelsDTO) (modelsInfo, error) {
	ms := dto.Models

	// get avatars
	v := map[string]userapp.UserDTO{}
	for i := range ms {
		v[ms[i].Owner] = userapp.UserDTO{}
	}

	accounts := make([]primitive.Account, len(v))
	i := 0
	for k := range v {
		accounts[i] = primitive.CreateAccount(k)
		userInfo, err := ctl.user.GetOrgOrUser(ctx, nil, accounts[i])
		if err != nil {
			return modelsInfo{}, err
		}
		v[k] = userInfo
		i++
	}

	// set avatars
	infos := make([]modelInfo, len(ms))
	for i := range ms {
		item := &ms[i]

		infos[i] = modelInfo{
			AvatarId:     v[item.Owner].AvatarId,
			OwnerType:    v[item.Owner].Type,
			Owner:        item.Owner,
			ModelSummary: item,
		}
	}

	return modelsInfo{
		Total:  dto.Total,
		Models: infos,
	}, nil
}

// @Summary  GetSpacesByModelId
// @Description  get spaces related to a model
// @Tags     ModelWeb
// @Param    id    path  string   true  "id of model" MaxLength(20)
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=[]spaceapp.SpaceModelDTO,msg=string,code=string}
// @Router   /v1/model/relation/{id}/space [get]
func (ctl *ModelWebController) GetSpacesByModelId(ctx *gin.Context) {
	modelId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	spaces, err := ctl.modelSpaceService.GetSpacesByModelId(ctx.Request.Context(), user, modelId)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	for _, space := range spaces {
		if avatar, err := ctl.user.GetUserAvatarId(ctx.Request.Context(), primitive.CreateAccount(space.Owner)); err != nil {
			space.AvatarId = ""
		} else {
			space.AvatarId = avatar.AvatarId
		}
	}

	commonctl.SendRespOfGet(ctx, &spaces)
}

// @Summary  disable model
// @Description  disable the model
// @Tags     ModelWeb
// @Param    id      path  string  true  "id of model" MaxLength(20)
// @Param    body  body      reqToDisableModel  true  "body of disable model"
// @Accept   json
// @Security Bearer
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Router   /v1/model/{id}/disable [put]
func (ctl *ModelWebController) disable(ctx *gin.Context) {
	middleware.SetAction(ctx, fmt.Sprintf("disable model of %s", ctx.Param("id")))

	req := reqToDisableModel{}
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	modelId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	action, err := ctl.appService.Disable(
		ctx.Request.Context(),
		ctl.userMiddleWare.GetUser(ctx),
		modelId,
		&cmd,
	)

	middleware.SetAction(ctx, fmt.Sprintf("%s, set %s", action, req.action()))

	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}
