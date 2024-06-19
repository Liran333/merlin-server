/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package controller

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/activity/app"
	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	modelapp "github.com/openmerlin/merlin-server/models/app"
	orgapp "github.com/openmerlin/merlin-server/organization/app"
	spaceapp "github.com/openmerlin/merlin-server/space/app"
	userapp "github.com/openmerlin/merlin-server/user/app"
)

// AddRouteForActivityRestfulController adds a router for the ActivityRestfulController with the given middleware.
func AddRouteForActivityRestfulController(
	r *gin.RouterGroup,
	s app.ActivityAppService,
	m middleware.UserMiddleWare,
	o orgapp.OrgService,
	u userapp.UserService,
	d modelapp.ModelAppService,
	p spaceapp.SpaceAppService,
	l middleware.OperationLog,
) {
	ctl := ActivityWebController{
		ActivityController: ActivityController{
			appService:     s,
			userMiddleWare: m,
			user:           u,
			org:            o,
			model:          d,
			space:          p,
		},
	}

	r.GET("/v1/user/activity", m.Optional, ctl.List)
	r.POST("/v1/like", m.Write, l.Write, ctl.Add)
	r.DELETE("/v1/like", m.Write, l.Write, ctl.Delete)
}

// ActivityRestfulController is a struct that holds the app service for model web operations.
type ActivityRestfulController struct {
	ActivityController
}

// @Summary  List
// @Description  get activities
// @Tags     ActivityRestful
// @Param    space query  string  false "filter by space" MaxLength(100)
// @Param    model query  string  false "filter by model" MaxLength(100)
// @Param    like  query  string  false "filter by like" MaxLength(100)
// @Accept   json
// @Security Bearer
// @Success  200  {object}  commonctl.ResponseData{data=activitiesInfo,msg=string,code=string}
// @Failure  400  {object}  commonctl.ResponseData{data=error,msg=string,code=string}
// @Router /v1/user/activity [get]
func (ctl *ActivityRestfulController) List(ctx *gin.Context) {
	// Bind query parameters to request struct
	var req reqToListUserActivities
	if err := ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}

	// Convert request to command
	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}

	// Get user from middleware
	user := ctl.userMiddleWare.GetUser(ctx)

	// Prepare list of names including the user's account name
	var list []primitive.Account

	// List activities based on the prepared list and command
	dto, err := ctl.appService.List(ctx.Request.Context(), user, list, &cmd)
	if err != nil {
		commonctl.SendError(ctx, err)
		return
	}

	if v, err := ctl.setAvatars(ctx, &dto); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, v)
	}
}

// setAvatars populates the avatar information for the activities in the provided ActivitysDTO.
func (ctl *ActivityRestfulController) setAvatars(ctx context.Context, dto *app.ActivitysDTO) (activitiesInfo, error) {
	ac := dto.Activities

	// get avatars
	v := map[string]bool{}
	for i := range ac {
		v[ac[i].Resource.Owner] = true
	}

	accounts := make([]primitive.Account, len(v))

	i := 0
	for k := range v {
		accounts[i] = primitive.CreateAccount(k)
		i++
	}

	avatars, err := ctl.user.GetUsersAvatarId(ctx, accounts)
	if err != nil {
		return activitiesInfo{}, err
	}

	// set avatars
	am := map[string]string{}
	for i := range avatars {
		item := &avatars[i]

		am[item.Name] = item.AvatarId
	}

	infos := make([]activityInfo, len(ac))
	for i := range ac {
		item := &ac[i]

		infos[i] = activityInfo{
			AvatarId:           am[item.Resource.Owner],
			ActivitySummaryDTO: item,
		}
	}

	return activitiesInfo{
		Total:      dto.Total,
		Activities: infos,
	}, nil
}
