/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package controller

import (
	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	userctl "github.com/openmerlin/merlin-server/user/controller"
	"github.com/openmerlin/merlin-server/utils"
)

func addRouteForSpaceSecretController(
	r *gin.RouterGroup,
	ctl *SpaceController,
	l middleware.OperationLog,
	sl middleware.SecurityLog,
	rl middleware.RateLimiter,
) {
	m := ctl.userMiddleWare

	r.POST(`/v1/space/:id/secret`, m.Write,
		userctl.CheckMail(ctl.userMiddleWare, ctl.user, sl), l.Write, rl.CheckLimit, ctl.CreateSecret)
	r.DELETE("/v1/space/:id/secret/:sid", m.Write,
		userctl.CheckMail(ctl.userMiddleWare, ctl.user, sl), l.Write, rl.CheckLimit, ctl.DeleteSecret)
	r.PUT("/v1/space/:id/secret/:sid", m.Write,
		userctl.CheckMail(ctl.userMiddleWare, ctl.user, sl), l.Write, rl.CheckLimit, ctl.UpdateSecret)
}

// @Summary  CreateSecret
// @Description  create space secret
// @Tags     Space
// @Param    body  body      reqToCreateSpaceSecret  true  "body of creating space secret"
// @Accept   json
// @Security Bearer
// @Success  201   {object}  commonctl.ResponseData{data=string,msg=string,code=string}
// @Router   /v1/space{id}/secret  [post]
func (ctl *SpaceController) CreateSecret(ctx *gin.Context) {
	req := reqToCreateSpaceSecret{}
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)
		return
	}

	defer utils.ClearStringMemory(*req.Value)

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	spaceId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}

	v, action, err := ctl.secretService.CreateSecret(user, spaceId, &cmd)

	middleware.SetAction(ctx, action)

	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, v)
	}
}

// @Summary  DeleteSecret
// @Description  delete space secret
// @Tags     Space
// @Param    id    path  string        true  "id of space" MaxLength(20)
// @Param    sid    path  string        true  "id of secret" MaxLength(20)
// @Accept   json
// @Security Bearer
// @Success  204
// @Router   /v1/space/{id}/secret/{sid} [delete]
func (ctl *SpaceController) DeleteSecret(ctx *gin.Context) {
	spaceId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}

	secretId, err := primitive.NewIdentity(ctx.Param("sid"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)
	action, err := ctl.secretService.DeleteSecret(user, spaceId, secretId)

	middleware.SetAction(ctx, action)

	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfDelete(ctx)
	}
}

// @Summary  UpdateSecret
// @Description  update space secret
// @Tags     Space
// @Param    id    path  string            true  "id of space" MaxLength(20)
// @Param    sid   path  string            true  "id of secret" MaxLength(20)
// @Param    body  body  reqToUpdateSpaceSecret  true  "body of updating space secret"
// @Accept   json
// @Security Bearer
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Router   /v1/space/{id}/secret/{vid} [put]
func (ctl *SpaceController) UpdateSecret(ctx *gin.Context) {
	req := reqToUpdateSpaceSecret{}
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)
		return
	}

	defer utils.ClearStringMemory(*req.Value)

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

	secretId, err := primitive.NewIdentity(ctx.Param("sid"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)
	action, err := ctl.secretService.UpdateSecret(
		user,
		spaceId, secretId, &cmd,
	)

	middleware.SetAction(ctx, action)

	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}
