/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides functionality for managing the application's controllers.
package controller

import (
	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/models/app"
)

// AddRouterForModelInternalController adds a router for the ModelInternalController with the given middleware.
func AddRouterForModelInternalController(
	r *gin.RouterGroup,
	s app.ModelInternalAppService,
	m middleware.UserMiddleWare,
) {
	ctl := ModelInternalController{
		appService: s,
	}

	r.GET("/v1/model/:id", m.Read, ctl.GetById)
	r.PUT("/v1/model/:id/label", m.Write, ctl.ResetLabel)
	r.PUT("/v1/model/:id", m.Write, ctl.Update)
	r.PUT("/v1/model/:id/use_in_openmind", m.Write, ctl.UpdateUseInOpenmind)
}

// ModelInternalController is a struct that holds the app service for model internal operations.
type ModelInternalController struct {
	appService app.ModelInternalAppService
}

// @Summary  ResetLabel
// @Description  reset label of model
// @Tags     ModelInternal
// @Param    id    path  string            true  "id of model"
// @Param    body  body  reqToResetLabel  true  "body"
// @Accept   json
// @Security Internal
// @Success  202  {object}  commonctl.ResponseData
// @Router   /v1/model/{id}/label [put]
func (ctl *ModelInternalController) ResetLabel(ctx *gin.Context) {
	req := reqToResetLabel{}

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	modelId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	cmd := req.toCmd()

	if err := ctl.appService.ResetLabels(modelId, &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  GetById
// @Description  get model info by id
// @Tags     ModelInternal
// @Param    id    path  string   true  "id of model"
// @Accept   json
// @Security Internal
// @Success  200  {object}  commonctl.ResponseData
// @Router   /v1/model/{id} [get]
func (ctl *ModelInternalController) GetById(ctx *gin.Context) {
	modelId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	data, err := ctl.appService.GetById(modelId)
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, data)
	}
}

// @Summary  Update model info
// @Description  update model info by id
// @Tags     ModelInternal
// @Param    id    path  string   true  "id of model"
// @Param    body  body  modelStatistics  true  "body of updating model info"
// @Accept   json
// @Security Internal
// @Success  202  {object}  commonctl.ResponseData
// @Router   /v1/model/{id} [put]
func (ctl *ModelInternalController) Update(ctx *gin.Context) {
	modelId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	var modelStatistics modelStatistics
	if err := ctx.BindJSON(&modelStatistics); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd := modelStatistics.toCmd()

	err = ctl.appService.UpdateStatistics(modelId, &cmd)
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  UpdateUseInOpenmind
// @Description  update space use in openmind info
// @Tags     ModelInternal
// @Param    id    path  string   true  "id of model"
// @Param    body  body  string   true  "use in openmind info"
// @Accept   json
// @Security Internal
// @Success  202  {object}  commonctl.ResponseData
// @Router   /v1/model/{id}/use_in_openmind [put]
func (ctl *ModelInternalController) UpdateUseInOpenmind(ctx *gin.Context) {
	spaceId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	var req useInOpenmind
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	envInfo := req.toCmd()

	err = ctl.appService.UpdateUseInOpenmind(spaceId, envInfo)
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}
