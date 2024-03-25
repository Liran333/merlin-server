/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-sdk/activityapp"
	"github.com/openmerlin/merlin-server/activity/app"
	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
)

// AddRouterForActivityInternalController adds a router for the ModelInternalController with the given middleware.
func AddRouterForActivityInternalController(
	r *gin.RouterGroup,
	s app.ActivityInternalAppService,
	m middleware.UserMiddleWare,
) {
	ctl := ActivityInternalController{
		appService: s,
	}

	r.POST("/v1/activity", m.Write, ctl.AddActivity)
	r.DELETE("/v1/activity", m.Write, ctl.DeleteActivity)
}

// ActivityInternalController is a struct that holds the app service for model internal operations.
type ActivityInternalController struct {
	appService app.ActivityInternalAppService
}

// @Summary  AddActivity
// @Description  add activities to DB
// @Tags     Activity
// @Accept   json
// @Security Bearer
// @Success  200  {object}  commonctl.ResponseData
// @Failure  400  {object}  commonctl.ResponseData
// @Router /v1/user/activity [post]
func (ctl *ActivityInternalController) AddActivity(ctx *gin.Context) {
	req := activityapp.ReqToCreateActivity{}

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := ConvertReqToCreateActivityToCmd(&req)

	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if err := ctl.appService.Create(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, nil)
	}
}

// @Summary  DeleteActivity
// @Description  delete all the record of an resource in the DB
// @Tags     Activity
// @Accept   json
// @Security Bearer
// @Success  200  {object}  commonctl.ResponseData
// @Failure  400  {object}  commonctl.ResponseData
// @Router /v1/user/activity [delete]
func (ctl *ActivityInternalController) DeleteActivity(ctx *gin.Context) {
	req := activityapp.ReqToDeleteActivity{}

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := ConvertReqToDeleteActivityToCmd(&req)

	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if err := ctl.appService.DeleteAll(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}
