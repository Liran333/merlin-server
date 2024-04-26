/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/space/app"
)

// AddRouterForSpaceInternalController adds routes to the given router group for the SpaceInternalController.
func AddRouterForSpaceInternalController(
	r *gin.RouterGroup,
	s app.SpaceInternalAppService,
	ms app.ModelSpaceAppService,
	m middleware.UserMiddleWare,
) {
	ctl := SpaceInternalController{
		appService:        s,
		modelSpaceService: ms,
	}

	r.GET("/v1/space/:id", m.Write, ctl.Get)
	r.PUT("/v1/space/:id/model", m.Write, ctl.UpdateSpaceModels)
	r.PUT("/v1/space/:id/local_cmd", m.Write, ctl.UpdateSpaceLocalCMD)
	r.PUT("/v1/space/:id/local_env_info", m.Write, ctl.UpdateSpaceLocalEnvInfo)
}

// SpaceInternalController is a struct that holds the necessary dependencies for handling space-related operations.
type SpaceInternalController struct {
	appService        app.SpaceInternalAppService
	modelSpaceService app.ModelSpaceAppService
}

// @Summary  Get
// @Description  get space
// @Tags     SpaceInternal
// @Param    id  path  string  true  "id of space"
// @Accept   json
// @Security Internal
// @Success  200  {object} commonctl.ResponseData
// @Router   /v1/space/{id} [get]
func (ctl *SpaceInternalController) Get(ctx *gin.Context) {
	spaceId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if dto, err := ctl.appService.GetById(spaceId); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, &dto)
	}
}

// @Summary  UpdateSpaceModels
// @Description  update space models relations
// @Tags     SpaceInternal
// @Param    id    path  string   true  "id of space"
// @Param    body  body  ModeIds  true  "body"
// @Accept   json
// @Security Internal
// @Success  202  {object}  commonctl.ResponseData
// @Router   /v1/space/{id}/model [put]
func (ctl *SpaceInternalController) UpdateSpaceModels(ctx *gin.Context) {
	spaceId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	_, err = ctl.appService.GetById(spaceId)
	if err != nil {
		commonctl.SendError(ctx, err)
	}

	var req ModeIds
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	modelsIndex := req.toCmd()

	err = ctl.modelSpaceService.UpdateRelation(spaceId, modelsIndex)
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  UpdateSpaceLocalCmd
// @Description  update space local cmd
// @Tags     SpaceInternal
// @Param    id    path  string   true  "id of space"
// @Param    body  body  string  true  "local cmd to reproduce the space"
// @Accept   json
// @Security Internal
// @Success  202  {object}  commonctl.ResponseData
// @Router   /v1/space/{id}/local_cmd [put]
func (ctl *SpaceInternalController) UpdateSpaceLocalCMD(ctx *gin.Context) {
	spaceId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	var req localCMD
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd := req.toCmd()

	err = ctl.appService.UpdateLocalCMD(spaceId, cmd)
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  UpdateSpaceLocalEnvInfo
// @Description  update space local env info
// @Tags     SpaceInternal
// @Param    id    path  string   true  "id of space"
// @Param    body  body  string   true  "local env info to update local space env info"
// @Accept   json
// @Security Internal
// @Success  202  {object}  commonctl.ResponseData
// @Router   /v1/space/{id}/local_env_info [put]
func (ctl *SpaceInternalController) UpdateSpaceLocalEnvInfo(ctx *gin.Context) {
	spaceId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	var req localEnvInfo
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	envInfo := req.toCmd()

	err = ctl.appService.UpdateEnvInfo(spaceId, envInfo)
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}
