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
	spacedomain "github.com/openmerlin/merlin-server/space/domain"
	"github.com/openmerlin/merlin-server/spaceapp/app"
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
	r.PUT(`/v1/space-app/build/started`, m.Write, ctl.NotifyBuildIsStarted)
	r.PUT(`/v1/space-app/build/done`, m.Write, ctl.NotifyBuildIsDone)
	r.PUT(`/v1/space-app/service/started`, m.Write, ctl.NotifyServiceIsStarted)
	r.PUT(`/v1/space-app/status`, m.Write, ctl.NotifyUpdateStatus)
	r.POST("/v1/space-app/:owner/:name/pause", m.Write, ctl.Pause)
}

// SpaceAppInternalController is a struct that holds the app service
// and provides methods for handling requests related to space apps.
type SpaceAppInternalController struct {
	appService 			app.SpaceappInternalAppService 
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

	if err := ctl.appService.Create(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, "successfully")
	}
}

// @Summary  NotifyBuildIsStarted
// @Description  notify space app building is started
// @Tags     SpaceApp
// @Param    body  body  reqToUpdateBuildInfo  true  "body"
// @Accept   json
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Security Internal
// @Router   /v1/space-app/build/started [put]
func (ctl *SpaceAppInternalController) NotifyBuildIsStarted(ctx *gin.Context) {
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

	if err := ctl.appService.NotifyBuildIsStarted(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  NotifyBuildIsDone
// @Description  notify space app build is done
// @Tags     SpaceApp
// @Param    body  body  reqToSetBuildIsDone  true  "body"
// @Accept   json
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Security Internal
// @Router   /v1/space-app/build/done [put]
func (ctl *SpaceAppInternalController) NotifyBuildIsDone(ctx *gin.Context) {
	req := reqToSetBuildIsDone{}

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if err := ctl.appService.NotifyBuildIsDone(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  NotifyServiceIsStarted
// @Description  notify space app service is started
// @Tags     SpaceApp
// @Param    body  body  reqToUpdateServiceInfo  true  "body"
// @Accept   json
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Security Internal
// @Router   /v1/space-app/service/started [put]
func (ctl *SpaceAppInternalController) NotifyServiceIsStarted(ctx *gin.Context) {
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

	if err := ctl.appService.NotifyServiceIsStarted(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  NotifyUpdateStatus
// @Description  notify space app status
// @Tags     SpaceApp
// @Param    body  body  reqToSetStatus  true  "body"
// @Accept   json
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Security Internal
// @Router   /v1/space-app/status [put]
func (ctl *SpaceAppInternalController) NotifyUpdateStatus(ctx *gin.Context) {
	req := reqToSetStatus{}

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if err := ctl.appService.NotifyUpdateStatus(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  Post
// @Description  stop space app
// @Tags     SpaceApp
// @Param    owner  path  string  true  "owner of space"
// @Param    name   path  string  true  "name of space"
// @Accept   json
// @Security Internal
// @Success  201   {object}  commonctl.ResponseData
// @Router   /v1/space-app/{owner}/{name}/pause [post]
func (ctl *SpaceAppInternalController) Pause(ctx *gin.Context) {
	index, err := ctl.parseIndex(ctx)
	if err != nil {
		return
	}

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

	if err := ctl.appService.PauseSpaceApp(&index, cmd.IsForce); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, "successfully")
	}
}

func (ctl *SpaceAppInternalController) parseIndex(ctx *gin.Context) (index spacedomain.SpaceIndex, err error) {
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
