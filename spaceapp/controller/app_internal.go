/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides the controllers for handling HTTP requests and managing the application's business logic.
package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/spaceapp/app"
	appprimitive "github.com/openmerlin/merlin-server/spaceapp/domain/primitive"
)

// AddRouteForSpaceappInternalController adds routes for SpaceAppInternalController to the given router group.
func AddRouteForSpaceappInternalController(
	r *gin.RouterGroup,
	s app.SpaceappInternalAppService,
	m middleware.UserMiddleWare,
) {

	ctl := SpaceAppInternalController{
		appService: s,
	}

	r.POST(`/v1/space-app`, m.Write, ctl.Create)

	r.PUT(`/v1/space-app/building`, m.Write, ctl.NotifySpaceAppBuilding)
	r.PUT(`/v1/space-app/starting`, m.Write, ctl.NotifySpaceAppStarting)
	r.PUT(`/v1/space-app/serving`, m.Write, ctl.NotifySpaceAppServing)
	r.PUT(`/v1/space-app/failed_status`, m.Write, ctl.NotifySpaceAppFailedStatus)

	r.POST(`/v1/space-app/pause`, m.Write, ctl.Pause)
	r.POST(`/v1/space-app/sleep`, m.Write, ctl.Sleep)
}

// SpaceAppInternalController is a struct that holds the app service
// and provides methods for handling requests related to space apps.
type SpaceAppInternalController struct {
	appService app.SpaceappInternalAppService
}

// @Summary  Create
// @Description  create space app
// @Tags     SpaceApp
// @Param    body  body  reqToCreateSpaceApp  true  "body of creating space app"
// @Accept   json
// @Success  201   {object}  commonctl.ResponseData
// @x-example {"data": "successfully"}
// @Security Internal
// @Router   /v1/space-app/ [post]
func (ctl *SpaceAppInternalController) Create(ctx *gin.Context) {
	req := reqToCreateSpaceApp{}
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if err := ctl.appService.Create(ctx.Request.Context(), &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, "successfully")
	}
}

// @Summary  NotifySpaceAppBuilding
// @Description  notify space app building is started
// @Tags     SpaceApp
// @Param    body  body  reqToUpdateBuildInfo  true  "body"
// @Accept   json
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Security Internal
// @Router   /v1/space-app/building [put]
func (ctl *SpaceAppInternalController) NotifySpaceAppBuilding(ctx *gin.Context) {
	req := reqToUpdateBuildInfo{}

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if err := ctl.appService.NotifyIsBuilding(ctx.Request.Context(), &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  NotifySpaceAppStarting
// @Description  notify space app build is starting
// @Tags     SpaceApp
// @Param    body  body  reqToNotifyStarting  true  "body"
// @Accept   json
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Security Internal
// @Router   /v1/space-app/starting [put]
func (ctl *SpaceAppInternalController) NotifySpaceAppStarting(ctx *gin.Context) {
	req := reqToNotifyStarting{}

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if err := ctl.appService.NotifyStarting(ctx.Request.Context(), &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  NotifySpaceAppServing
// @Description  notify space app service is started
// @Tags     SpaceApp
// @Param    body  body  reqToUpdateServiceInfo  true  "body"
// @Accept   json
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Security Internal
// @Router   /v1/space-app/serving [put]
func (ctl *SpaceAppInternalController) NotifySpaceAppServing(ctx *gin.Context) {
	req := reqToUpdateServiceInfo{}

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if err := ctl.appService.NotifyIsServing(ctx.Request.Context(), &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  NotifySpaceAppFailedStatus
// @Description  notify space app failed status
// @Tags     SpaceApp
// @Param    body  body  reqToFailedStatus  true  "body"
// @Accept   json
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Security Internal
// @Router   /v1/space-app/failed_status [put]
func (ctl *SpaceAppInternalController) NotifySpaceAppFailedStatus(ctx *gin.Context) {
	req := reqToFailedStatus{}

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	switch cmd.Status {
	case appprimitive.AppStatusBuildFailed:
		if err := ctl.appService.NotifyIsBuildFailed(ctx.Request.Context(), &cmd); err != nil {
			commonctl.SendError(ctx, err)
			return
		}
	case appprimitive.AppStatusStartFailed:
		if err := ctl.appService.NotifyIsStartFailed(ctx.Request.Context(), &cmd); err != nil {
			commonctl.SendError(ctx, err)
			return
		}
	case appprimitive.AppStatusRestartFailed:
		if err := ctl.appService.NotifyIsRestartFailed(ctx.Request.Context(), &cmd); err != nil {
			commonctl.SendError(ctx, err)
			return
		}
	case appprimitive.AppStatusResumeFailed:
		if err := ctl.appService.NotifyIsResumeFailed(ctx.Request.Context(), &cmd); err != nil {
			commonctl.SendError(ctx, err)
			return
		}
	default:
		e := fmt.Errorf("old status not %s, can not set", cmd.Status.AppStatus())
		err = allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
		commonctl.SendError(ctx, err)
		return
	}
	commonctl.SendRespOfPost(ctx, "successfully")
}

// @Summary  Post
// @Description  pause space app
// @Tags     SpaceApp
// @Param    owner  path  string  true  "owner of space"
// @Param    name   path  string  true  "name of space"
// @Accept   json
// @Security Internal
// @Success  201   {object}  commonctl.ResponseData
// @Router   /v1/space-app/pause [post]
func (ctl *SpaceAppInternalController) Pause(ctx *gin.Context) {

	req := reqToPauseSpaceApp{}

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)
		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}
	if cmd.IsForce {
		if err := ctl.appService.ForcePauseSpaceApp(ctx.Request.Context(), cmd.SpaceId); err != nil {
			commonctl.SendError(ctx, err)
			return
		}
	} else {
		if err := ctl.appService.PauseSpaceApp(ctx.Request.Context(), cmd.SpaceId); err != nil {
			commonctl.SendError(ctx, err)
			return
		}
	}
	commonctl.SendRespOfPost(ctx, "successfully")
}

// @Summary  Post
// @Description  sleep space app
// @Tags     SpaceApp
// @Param    owner  path  string  true  "owner of space"
// @Param    name   path  string  true  "name of space"
// @Accept   json
// @Security Internal
// @Success  201   {object}  commonctl.ResponseData
// @Router   /v1/space-app/sleep [post]
func (ctl *SpaceAppInternalController) Sleep(ctx *gin.Context) {

	req := reqToSleepSpaceApp{}

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)
		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}
	if err := ctl.appService.SleepSpaceApp(ctx.Request.Context(), &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, "successfully")
	}
}
