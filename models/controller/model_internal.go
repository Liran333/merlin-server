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
	spaceapp "github.com/openmerlin/merlin-server/space/app"
)

// AddRouterForModelInternalController adds a router for the ModelInternalController with the given middleware.
func AddRouterForModelInternalController(
	r *gin.RouterGroup,
	s app.ModelInternalAppService,
	ms spaceapp.ModelSpaceAppService,
	m middleware.UserMiddleWare,
) {
	ctl := ModelInternalController{
		appService:        s,
		modelSpaceService: ms,
	}

	r.GET("/v1/model/:id", m.Read, ctl.GetById)
	r.PUT("/v1/model/:id/label", m.Write, ctl.ResetLabel)
	r.PUT("/v1/model/:id", m.Write, ctl.Update)

	r.PUT("/v1/model/:id/use_in_openmind", m.Write, ctl.UpdateUseInOpenmind)
	r.GET("/v1/model/relation/:id/space", m.Read, ctl.GetSpacesByModelId)

	r.PUT("/v1/model/deploy/:owner/:name", m.Write, ctl.Deploy)
}

// ModelInternalController is a struct that holds the app service for model internal operations.
type ModelInternalController struct {
	ModelController
	appService        app.ModelInternalAppService
	modelSpaceService spaceapp.ModelSpaceAppService
}

// @Summary  ResetLabel
// @Description  reset label of model
// @Tags     ModelInternal
// @Param    id    path  string            true  "id of model" MaxLength(20)
// @Param    body  body  reqToResetLabel  true  "body"
// @Accept   json
// @Security Internal
// @Success  202  {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
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
// @Param    id    path  string   true  "id of model" MaxLength(20)
// @Accept   json
// @Security Internal
// @Success  200  {object}  commonctl.ResponseData{data=app.ModelDTO,msg=string,code=string}
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
// @Param    id    path  string   true  "id of model" MaxLength(20)
// @Param    body  body  modelStatistics  true  "body of updating model info"
// @Accept   json
// @Security Internal
// @Success  202  {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
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
// @Param    id    path  string   true  "id of model" MaxLength(20)
// @Param    body  body  string   true  "use in openmind info"
// @Accept   json
// @Security Internal
// @Success  202  {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
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

// @Summary  GetSpacesByModelId
// @Description  get all spaces related to a model, including those that have been disabled.
// @Tags     ModelInternal
// @Param    id    path  string   true  "id of model" MaxLength(20)
// @Accept   json
// @Security Internal
// @Success  200  {object}  commonctl.ResponseData{data=app.SpaceIdModelDTO,msg=string,code=string}
// @Router   /v1/model/relation/{id}/space [get]
func (ctl *ModelInternalController) GetSpacesByModelId(ctx *gin.Context) {
	modelId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	spaces, err := ctl.modelSpaceService.GetSpaceIdsByModelId(modelId)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	commonctl.SendRespOfGet(ctx, &spaces)
}

// @Summary  Update deploy
// @Description  update deploy info of model
// @Tags     ModelInternal
// @Param    owner  path  string  true  "owner of model" MaxLength(40)
// @Param    name   path  string  true  "name of model" MaxLength(100)
// @Param    body  body  reqToSaveModelDeploy   true  "deploy info"
// @Accept   json
// @Security Internal
// @Success  202  {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Router   /v1/model/deploy/{owner}/{name} [put]
func (ctl *ModelInternalController) Deploy(ctx *gin.Context) {
	index, err := ctl.parseIndex(ctx)
	if err != nil {
		return
	}

	var req reqToSaveModelDeploy
	if err = ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	err = ctl.appService.SaveDeploy(index, req)
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}
