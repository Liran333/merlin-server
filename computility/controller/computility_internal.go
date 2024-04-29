/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package controller

import (
	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/computility/app"
)

// AddRouterForComputilityInternalController adds routes to the given router group for the ComputilityInternalController.
func AddRouterForComputilityInternalController(
	r *gin.RouterGroup,
	s app.ComputilityInternalAppService,
	m middleware.UserMiddleWare,
) {
	ctl := ComputilityInternalController{
		appService: s,
	}

	r.POST("/v1/computility/account", m.Write, ctl.ComputilityUserJoin)
	r.PUT("/v1/computility/account/remove", m.Write, ctl.ComputilityUserRemove)
	r.POST("/v1/computility/org/delete", m.Write, ctl.ComputilityOrgDelete)
}

// ComputilityInternalController is a struct that holds the necessary dependencies for handling computility-related operations.
type ComputilityInternalController struct {
	appService app.ComputilityInternalAppService
}

// @Summary  ComputilityUserJoin
// @Description  user joined computility org
// @Tags     ComputilityInternal
// @Param    body  body  reqToUserOrgOperate  true  "body"
// @Accept   json
// @Security Internal
// @Success  201   {object} commonctl.ResponseData{data=nil,msg=string,code=string}
// @Failure  400   {object} commonctl.ResponseData{data=error,msg=string,code=string}
// @Router   /v1/computility/account [post]
func (ctl *ComputilityInternalController) ComputilityUserJoin(ctx *gin.Context) {
	req := reqToUserOrgOperate{}

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	err = ctl.appService.UserJoin(cmd)
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, nil)
	}
}

// @Summary  ComputilityUserRemove
// @Description  user removed from computility org
// @Tags     ComputilityInternal
// @Param    body  body  reqToUserOrgOperate  true  "body"
// @Accept   json
// @Success  202   {object}  commonctl.ResponseData{data=app.AccountRecordlDTO,msg=string,code=string}
// @Failure  400   {object}  commonctl.ResponseData{data=error,msg=string,code=string}
// @Security Internal
// @Router   /v1/computility/account/remove [put]
func (ctl *ComputilityInternalController) ComputilityUserRemove(ctx *gin.Context) {
	req := reqToUserOrgOperate{}

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	r, err := ctl.appService.UserRemove(cmd)
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, r)
	}
}

// @Summary  ComputilityOrgDelete
// @Description  delete computility org
// @Tags     ComputilityInternal
// @Param    body  body  reqToOrgDelete  true  "body"
// @Accept   json
// @Success  204   {object}  commonctl.ResponseData{data=[]app.AccountRecordlDTO,msg=string,code=string}
// @Failure  400   {object}  commonctl.ResponseData{data=error,msg=string,code=string}
// @Security Internal
// @Router   /v1/computility/org/delete [post]
func (ctl *ComputilityInternalController) ComputilityOrgDelete(ctx *gin.Context) {
	req := reqToOrgDelete{}

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	r, err := ctl.appService.OrgDelete(cmd)
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, r)
	}
}
