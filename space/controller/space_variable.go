/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	userctl "github.com/openmerlin/merlin-server/user/controller"
)

func addRouteForSpaceVariableController(
	r *gin.RouterGroup,
	ctl *SpaceController,
	l middleware.OperationLog,
	sl middleware.SecurityLog,
	rl middleware.RateLimiter,
) {
	m := ctl.userMiddleWare

	r.POST(`/v1/space/:id/variable`, m.Write,
		userctl.CheckMail(ctl.userMiddleWare, ctl.user, sl), l.Write, rl.CheckLimit, ctl.CreateVariable)
	r.DELETE("/v1/space/:id/variable/:vid", m.Write,
		userctl.CheckMail(ctl.userMiddleWare, ctl.user, sl), l.Write, rl.CheckLimit, ctl.DeleteVariable)
	r.PUT("/v1/space/:id/variable/:vid", m.Write,
		userctl.CheckMail(ctl.userMiddleWare, ctl.user, sl), l.Write, rl.CheckLimit, ctl.UpdateVariable)
	r.GET("/v1/space/:owner/:name/variable-secret", m.Read,
		userctl.CheckMail(ctl.userMiddleWare, ctl.user, sl), l.Write, rl.CheckLimit, ctl.GetVariableSecret)
}

// @Summary  CreateVariable
// @Description  create space variable
// @Tags     Space
// @Param    body  body      reqToCreateSpaceVariable  true  "body of creating space variable"
// @Accept   json
// @Security Bearer
// @Success  201   {object}  commonctl.ResponseData{data=string,msg=string,code=string}
// @Router   /v1/space/{id}/variable [post]
func (ctl *SpaceController) CreateVariable(ctx *gin.Context) {
	req := reqToCreateSpaceVariable{}
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)
		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}

	spaceId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)
	v, action, err := ctl.variableService.CreateVariable(user, spaceId, &cmd)

	middleware.SetAction(ctx, action)

	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, v)
	}
}

// @Summary  DeleteVariable
// @Description  delete space variable
// @Tags     Space
// @Param    id    path  string        true  "id of space" MaxLength(20)
// @Param    vid    path  string        true  "id of variable" MaxLength(20)
// @Accept   json
// @Security Bearer
// @Success  204
// @Router   /v1/space/{id}/variable/{vid} [delete]
func (ctl *SpaceController) DeleteVariable(ctx *gin.Context) {
	spaceId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}

	variableId, err := primitive.NewIdentity(ctx.Param("vid"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)
	action, err := ctl.variableService.DeleteVariable(user, spaceId, variableId)

	middleware.SetAction(ctx, action)

	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfDelete(ctx)
	}
}

// @Summary  UpdateVariable
// @Description  update space variable
// @Tags     Space
// @Param    id     path  string            true  "id of space" MaxLength(20)
// @Param    vid    path  string            true  "id of variable" MaxLength(20)
// @Param    body   body  reqToUpdateSpaceVariable  true  "body of updating space variable"
// @Accept   json
// @Security Bearer
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Router   /v1/space/{id}/variable/{vid} [put]
func (ctl *SpaceController) UpdateVariable(ctx *gin.Context) {
	req := reqToUpdateSpaceVariable{}
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)
		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}

	spaceId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}

	variableId, err := primitive.NewIdentity(ctx.Param("vid"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)
	action, err := ctl.variableService.UpdateVariable(
		user, spaceId, variableId, &cmd,
	)

	middleware.SetAction(ctx, fmt.Sprintf("%s, set %s", action, req.action()))

	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  List
// @Description  list space variable secret
// @Tags     SpaceWeb
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=app.SpaceVariableSecretDTO,msg=string,code=string}
// @Router   /v1/space/:owner/:name/variable-secret [get]
func (ctl *SpaceController) GetVariableSecret(ctx *gin.Context) {
	index, err := ctl.parseIndex(ctx)
	if err != nil {
		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)
	space, err := ctl.appService.GetByName(user, &index)
	if err != nil {
		commonctl.SendError(ctx, err)
		return
	}

	dto, err := ctl.variableService.ListVariableSecret(space.Id)
	if err != nil {
		commonctl.SendError(ctx, err)
		return
	} else {
		commonctl.SendRespOfGet(ctx, &dto)
	}
}
