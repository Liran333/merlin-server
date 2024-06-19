/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides functionality for managing the application's controllers.
package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	activityapp "github.com/openmerlin/merlin-server/activity/app"
	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/models/app"
	"github.com/openmerlin/merlin-server/models/domain"
	userapp "github.com/openmerlin/merlin-server/user/app"
	userctl "github.com/openmerlin/merlin-server/user/controller"
)

func addRouteForModelController(
	r *gin.RouterGroup,
	ctl *ModelController,
	opLog middleware.OperationLog,
	sl middleware.SecurityLog,
) {
	m := ctl.userMiddleWare

	r.POST(`/v1/model`, m.Write, userctl.CheckMail(ctl.userMiddleWare, ctl.user, sl), opLog.Write, ctl.Create)
	r.DELETE("/v1/model/:id", m.Write, userctl.CheckMail(ctl.userMiddleWare, ctl.user, sl), opLog.Write, ctl.Delete)
	r.PUT("/v1/model/:id", m.Write, userctl.CheckMail(ctl.userMiddleWare, ctl.user, sl), opLog.Write, ctl.Update)
}

// ModelController is a controller for handling model-related requests.
type ModelController struct {
	user           userapp.UserService
	appService     app.ModelAppService
	userMiddleWare middleware.UserMiddleWare
	activity       activityapp.ActivityAppService
}

// @Summary  Create
// @Description  create model
// @Tags     Model
// @Param    body  body      reqToCreateModel  true  "body of creating model"
// @Accept   json
// @Security Bearer
// @Success  201   {object}  commonctl.ResponseData{data=string,msg=string,code=string}
// @Router   /v1/model [post]
func (ctl *ModelController) Create(ctx *gin.Context) {
	middleware.SetAction(ctx, "create model")

	req := reqToCreateModel{}
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	middleware.SetAction(ctx, req.action())

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	if v, err := ctl.appService.Create(ctx, user, &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, v)
	}
}

// @Summary  Delete
// @Description  delete model
// @Tags     Model
// @Param    id    path  string        true  "id of model" MaxLength(20)
// @Accept   json
// @Security Bearer
// @Success  204
// @Router   /v1/model/{id} [delete]
func (ctl *ModelController) Delete(ctx *gin.Context) {
	middleware.SetAction(ctx, fmt.Sprintf("delete model of %s", ctx.Param("id")))

	user := ctl.userMiddleWare.GetUser(ctx)

	modelId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	action, err := ctl.appService.Delete(ctx.Request.Context(), user, modelId)

	middleware.SetAction(ctx, action)

	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfDelete(ctx)
	}
}

// @Summary  Update
// @Description  update model
// @Tags     Model
// @Param    id    path  string            true  "id of model" MaxLength(20)
// @Param    body  body  reqToUpdateModel  true  "body of updating model"
// @Accept   json
// @Security Bearer
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Router   /v1/model/{id} [put]
func (ctl *ModelController) Update(ctx *gin.Context) {
	middleware.SetAction(ctx, fmt.Sprintf("update model of %s", ctx.Param("id")))

	req := reqToUpdateModel{}
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

	action, err := ctl.appService.Update(
		ctx.Request.Context(), ctl.userMiddleWare.GetUser(ctx),
		modelId, &cmd,
	)

	middleware.SetAction(ctx, fmt.Sprintf("%s, set %s", action, req.action()))

	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

func (ctl *ModelController) parseIndex(ctx *gin.Context) (index domain.ModelIndex, err error) {
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
