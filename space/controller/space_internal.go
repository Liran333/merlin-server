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
