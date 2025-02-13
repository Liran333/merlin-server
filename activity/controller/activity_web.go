/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// nolint:typecheck
package controller

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/openmerlin/merlin-sdk/activityapp"

	"github.com/openmerlin/merlin-server/activity/app"
	"github.com/openmerlin/merlin-server/common/controller"
	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	datasetapp "github.com/openmerlin/merlin-server/datasets/app"
	modelapp "github.com/openmerlin/merlin-server/models/app"
	orgapp "github.com/openmerlin/merlin-server/organization/app"
	spaceapp "github.com/openmerlin/merlin-server/space/app"
	userapp "github.com/openmerlin/merlin-server/user/app"
)

// AddRouteForActivityWebController adds a router for the ActivityWebController with the given middleware.
func AddRouteForActivityWebController(
	r *gin.RouterGroup,
	s app.ActivityAppService,
	m middleware.UserMiddleWare,
	o orgapp.OrgService,
	u userapp.UserService,
	d modelapp.ModelAppService,
	ds datasetapp.DatasetAppService,
	p spaceapp.SpaceAppService,
	rl middleware.RateLimiter,
	l middleware.OperationLog,
) {
	ctl := ActivityWebController{
		ActivityController: ActivityController{
			appService:     s,
			userMiddleWare: m,
			user:           u,
			org:            o,
			model:          d,
			dataset:        ds,
			space:          p,
		},
	}

	r.GET("/v1/user/activity", m.Read, rl.CheckLimit, ctl.List)
	r.POST("/v1/like", m.Write, l.Write, ctl.Add)
	r.DELETE("/v1/like", m.Write, l.Write, ctl.Delete)
}

// toCmd converts a reqToListUserActivities object to an app.CmdToListActivities command.
func (req *reqToListUserActivities) toCmd() (cmd app.CmdToListActivities, err error) {
	cmd.Count = req.Count
	cmd.Model = req.Model
	cmd.Dataset = req.Dataset
	cmd.Space = req.Space
	cmd.Like = req.Like
	if v := req.CountPerPage; v <= 0 || v > config.MaxCountPerPage {
		cmd.CountPerPage = config.MaxCountPerPage
	} else {
		cmd.CountPerPage = v
	}

	if v := req.PageNum; v <= 0 {
		cmd.PageNum = firstPage
	} else {
		if v > (math.MaxInt / cmd.CountPerPage) {
			err = errors.New("invalid page num")

			return
		}
		cmd.PageNum = v
	}

	return
}

// ActivityWebController is a struct that holds the app service for model web operations.
type ActivityWebController struct {
	ActivityController
}

// reqToListUserModels
type reqToListUserActivities struct {
	Model   string `form:"model"`
	Dataset string `form:"dataset"`
	Space   string `form:"space"`
	Like    string `form:"like"`
	controller.CommonListRequest
}

// @Summary  List
// @Description  get activities
// @Tags     ActivityWeb
// @Param    space     query  string  false "filter by space" MaxLength(100)
// @Param    model  query  string  false "filter by model" MaxLength(100)
// @Param    dataset  query  string  false "filter by dataset" MaxLength(100)
// @Param    like  query  string  false "filter by like" MaxLength(100)
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=app.ActivitysDTO,msg=string,code=string}
// @Failure  400  {object}  commonctl.ResponseData{data=error,msg=string,code=string}
// @Router /v1/user/activity [get]
func (ctl *ActivityWebController) List(ctx *gin.Context) {
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

	var list []primitive.Account

	list = append(list, user)

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

// @Summary  Add
// @Description  add a like record in the activity table
// @Tags     ActivityWeb
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Failure  400  {object}  commonctl.ResponseData{data=error,msg=string,code=string}
// @Router /web/v1/like [post]
func (ctl *ActivityWebController) Add(ctx *gin.Context) {
	middleware.SetAction(ctx, "start to add a like")

	var req activityapp.ReqToCreateActivity

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)
		return
	}

	actionDescription := fmt.Sprintf("add a like to a %s, id: %v", req.ResourceType, req.ResourceId)
	middleware.SetAction(ctx, actionDescription)

	user := ctl.userMiddleWare.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	req.Owner = user.Account()

	cmd, err := ConvertReqToCreateActivityToCmd(&req)

	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}

	liked, err := ctl.appService.HasLike(user, cmd.Resource.Index)
	if err != nil {
		return
	}
	if !liked {
		if req.ResourceType == typeModel {
			err = ctl.model.AddLike(cmd.Resource.Index)
		} else if req.ResourceType == typeDataset {
			err = ctl.dataset.AddLike(cmd.Resource.Index)
		} else {
			err = ctl.space.AddLike(cmd.Resource.Index)
		}
		// Check for errors from AddLike operation
		if err != nil {
			commonctl.SendError(ctx, err)
			return
		}
	}

	if err := ctl.appService.Create(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  Delete
// @Description  Delete a like record in the activity table
// @Tags     ActivityWeb
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData
// @Failure  400  {object}  commonctl.ResponseData{data=error,msg=string,code=string}
// @Router /web/v1/like [delete]
func (ctl *ActivityWebController) Delete(ctx *gin.Context) {
	middleware.SetAction(ctx, "start to delete a like")

	var req activityapp.ReqToDeleteActivity

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)
		return
	}

	actionDescription := fmt.Sprintf("cancel a like to a %s, id: %v", req.ResourceType, req.ResourceId)
	middleware.SetAction(ctx, actionDescription)

	user := ctl.userMiddleWare.GetUserAndExitIfFailed(ctx)
	if user == nil {
		return
	}

	cmd, err := ConvertReqToDeleteActivityToCmd(user, &req)

	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}

	liked, err := ctl.appService.HasLike(user, cmd.Resource.Index)
	if err != nil {
		return
	}
	if liked {
		if req.ResourceType == typeModel {
			err = ctl.model.DeleteLike(cmd.Resource.Index)
		} else if req.ResourceType == typeDataset {
			err = ctl.dataset.DeleteLike(cmd.Resource.Index)
		} else {
			err = ctl.space.DeleteLike(cmd.Resource.Index)
		}

		// Check for errors from DeleteLike operation
		if err != nil {
			commonctl.SendError(ctx, err)
			return
		}
	}

	if err := ctl.appService.Delete(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfDelete(ctx)
	}
}

// setAvatars sets avatars for the activities in an ActivitysDTO object.
func (ctl *ActivityWebController) setAvatars(ctx context.Context, dto *app.ActivitysDTO) (activitiesInfo, error) {
	as := dto.Activities

	// get avatars
	v := map[string]userapp.UserDTO{}
	for i := range as {
		v[as[i].Resource.Owner] = userapp.UserDTO{}
	}

	accounts := make([]primitive.Account, len(v))

	i := 0
	for k := range v {
		accounts[i] = primitive.CreateAccount(k)
		userInfo, err := ctl.user.GetOrgOrUser(ctx, nil, accounts[i])
		if err != nil {
			return activitiesInfo{}, err
		}
		v[k] = userInfo
		i++
	}

	infos := make([]activityInfo, len(as))
	for i := range as {
		item := &as[i]

		infos[i] = activityInfo{
			AvatarId:           v[item.Resource.Owner].AvatarId,
			OwnerType:          v[item.Resource.Owner].Type,
			ActivitySummaryDTO: item,
		}
	}

	return activitiesInfo{
		Total:      dto.Total,
		Activities: infos,
	}, nil
}
