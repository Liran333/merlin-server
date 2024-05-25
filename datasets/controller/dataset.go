/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides functionality for managing the application's controllers.
package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"

	activityapp "github.com/openmerlin/merlin-server/activity/app"
	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/datasets/app"
	"github.com/openmerlin/merlin-server/datasets/domain"
	userapp "github.com/openmerlin/merlin-server/user/app"
	userctl "github.com/openmerlin/merlin-server/user/controller"
)

func addRouteForDatasetController(
	r *gin.RouterGroup,
	ctl *DatasetController,
	opLog middleware.OperationLog,
	sl middleware.SecurityLog,
) {
	m := ctl.userMiddleWare

	r.POST(`/v1/dataset`, m.Write, userctl.CheckMail(ctl.userMiddleWare, ctl.user, sl), opLog.Write, ctl.Create)
	r.DELETE("/v1/dataset/:id", m.Write, userctl.CheckMail(ctl.userMiddleWare, ctl.user, sl), opLog.Write, ctl.Delete)
	r.PUT("/v1/dataset/:id", m.Write, userctl.CheckMail(ctl.userMiddleWare, ctl.user, sl), opLog.Write, ctl.Update)
}

// DatasetController is a controller for handling dataset-related requests.
type DatasetController struct {
	user           userapp.UserService
	appService     app.DatasetAppService
	userMiddleWare middleware.UserMiddleWare
	activity       activityapp.ActivityAppService
}

// @Summary  Create
// @Description  create dataset
// @Tags     Dataset
// @Param    body  body      reqToCreateDataset  true  "body of creating dataset"
// @Accept   json
// @Security Bearer
// @Success  201   {object}  commonctl.ResponseData{data=string,msg=string,code=string}
// @Router   /v1/dataset [post]
func (ctl *DatasetController) Create(ctx *gin.Context) {
	middleware.SetAction(ctx, "create dataset")

	req := reqToCreateDataset{}
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, xerrors.Errorf("failed to parse req, %w", err))

		return
	}

	middleware.SetAction(ctx, req.action())

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to convert req to cmd, %w", err))

		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	if v, err := ctl.appService.Create(user, &cmd); err != nil {
		commonctl.SendError(ctx, xerrors.Errorf("create dataset failed, err:%w", err))
	} else {
		commonctl.SendRespOfPost(ctx, v)
	}
}

// @Summary  Delete
// @Description  delete dataset
// @Tags     Dataset
// @Param    id    path  string        true  "id of dataset" MaxLength(20)
// @Accept   json
// @Security Bearer
// @Success  204
// @Router   /v1/dataset/{id} [delete]
func (ctl *DatasetController) Delete(ctx *gin.Context) {
	middleware.SetAction(ctx, fmt.Sprintf("delete dataset of %s", ctx.Param("id")))

	user := ctl.userMiddleWare.GetUser(ctx)

	datasetId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("%w", err))

		return
	}

	action, err := ctl.appService.Delete(user, datasetId)

	middleware.SetAction(ctx, action)

	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfDelete(ctx)
	}
}

// @Summary  Update
// @Description  update dataset
// @Tags     Dataset
// @Param    id    path  string            true  "id of dataset" MaxLength(20)
// @Param    body  body  reqToUpdateDataset  true  "body of updating dataset"
// @Accept   json
// @Security Bearer
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Router   /v1/dataset/{id} [put]
func (ctl *DatasetController) Update(ctx *gin.Context) {
	middleware.SetAction(ctx, fmt.Sprintf("update dataset of %s", ctx.Param("id")))

	req := reqToUpdateDataset{}
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, xerrors.Errorf("failed to parse req, %w", err))

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to convert req to cmd, %w", err))

		return
	}

	datasetId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("%w", err))

		return
	}

	action, err := ctl.appService.Update(
		ctl.userMiddleWare.GetUser(ctx),
		datasetId, &cmd,
	)

	middleware.SetAction(ctx, fmt.Sprintf("%s, set %s", action, req.action()))

	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

func (ctl *DatasetController) parseIndex(ctx *gin.Context) (index domain.DatasetIndex, err error) {
	index.Owner, err = primitive.NewAccount(ctx.Param("owner"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("%w", err))

		return
	}

	index.Name, err = primitive.NewMSDName(ctx.Param("name"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("%w", err))
	}

	return
}
