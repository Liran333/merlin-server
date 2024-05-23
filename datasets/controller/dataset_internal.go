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
	"github.com/openmerlin/merlin-server/datasets/app"
)

// AddRouterForDatasetInternalController adds a router for the DatasetInternalController with the given middleware.
func AddRouterForDatasetInternalController(
	r *gin.RouterGroup,
	s app.DatasetInternalAppService,
	m middleware.UserMiddleWare,
) {
	ctl := DatasetInternalController{
		appService: s,
	}

	r.GET("/v1/dataset/:id", m.Read, ctl.GetById)
	r.PUT("/v1/dataset/:id/label", m.Write, ctl.ResetLabel)
	r.PUT("/v1/dataset/:id", m.Write, ctl.Update)
}

// DatasetInternalController is a struct that holds the app service for dataset internal operations.
type DatasetInternalController struct {
	appService app.DatasetInternalAppService
}

// @Summary  ResetLabel
// @Description  reset label of datasets
// @Tags     DatasetInternal
// @Param    id    path  string            true  "id of dataset" MaxLength(20)
// @Param    body  body  reqToResetDatasetLabel  true  "body"
// @Accept   json
// @Security Internal
// @Success  202  {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Router   /v1/dataset/{id}/label [put]
func (ctl *DatasetInternalController) ResetLabel(ctx *gin.Context) {
	req := reqToResetDatasetLabel{}

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	datasetId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	cmd := req.toCmd()

	if err := ctl.appService.ResetLabels(datasetId, &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  GetById
// @Description  get dataset info by id
// @Tags     DatasetInternal
// @Param    id    path  string   true  "id of dataset" MaxLength(20)
// @Accept   json
// @Security Internal
// @Success  200  {object}  commonctl.ResponseData{data=app.DatasetDTO,msg=string,code=string}
// @Router   /v1/dataset/{id} [get]
func (ctl *DatasetInternalController) GetById(ctx *gin.Context) {
	datasetId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	data, err := ctl.appService.GetById(datasetId)
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, data)
	}
}

// @Summary  Update dataset info
// @Description  update dataset info by id
// @Tags     DatasetInternal
// @Param    id    path  string   true  "id of dataset" MaxLength(20)
// @Param    body  body  datasetStatistics  true  "body of updating dataset info"
// @Accept   json
// @Security Internal
// @Success  202  {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Router   /v1/dataset/{id} [put]
func (ctl *DatasetInternalController) Update(ctx *gin.Context) {
	datasetId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	var datasetStatistics datasetStatistics
	if err := ctx.BindJSON(&datasetStatistics); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd := datasetStatistics.toCmd()

	err = ctl.appService.UpdateStatistics(datasetId, &cmd)
	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}
