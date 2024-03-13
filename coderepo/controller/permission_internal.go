/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/coderepo/domain/resourceadapter"
	commonapp "github.com/openmerlin/merlin-server/common/app"
	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
)

//nolint:golint,unused
var errorNoPermission = allerror.NewNoPermission("no permission")

// AddRouteForCodeRepoPermissionInternalController adds routes for CodeRepoPermissionInternalController.
func AddRouteForCodeRepoPermissionInternalController(
	r *gin.RouterGroup,
	s commonapp.ResourcePermissionAppService,
	a resourceadapter.ResourceAdapter,
	m middleware.UserMiddleWare,
) {

	ctl := PermissionInternalController{
		ps:   s,
		repo: a,
	}

	r.POST(`/v1/coderepo/permission/update`, m.Write, ctl.Update)
	r.POST(`/v1/coderepo/permission/read`, m.Write, ctl.Read)
}

// PermissionInternalController is a struct that holds the necessary services
// and adapters for handling permission-related operations.
type PermissionInternalController struct {
	ps   commonapp.ResourcePermissionAppService
	repo resourceadapter.ResourceAdapter
}

// @Summary  Update
// @Description  check if can create/update/delete repo's sub-resource not the repo itsself
// @Tags     Permission
// @Param    body  body  reqToCheckPermission  true  "body of request"
// @Accept   json
// @Success  201   {object}  commonctl.ResponseData
// @Security Internal
// @Router   /v1/coderepo/permission/update [post]
func (ctl *PermissionInternalController) Update(ctx *gin.Context) {
	user, r, err := ctl.parse(ctx)
	if err != nil {
		return
	}

	if err := ctl.ps.CanUpdate(user, r); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, "successfully")
	}
}

// @Summary  Read
// @Description  check if can read repo's sub-resource not the repo itsself
// @Tags     Permission
// @Param    body  body  reqToCheckPermission  true  "body of request"
// @Accept   json
// @Success  201   {object}  commonctl.ResponseData
// @Security Internal
// @Router   /v1/coderepo/permission/read [post]
func (ctl *PermissionInternalController) Read(ctx *gin.Context) {
	user, r, err := ctl.parse(ctx)
	if err != nil {
		return
	}

	if err := ctl.ps.CanRead(user, r); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, "successfully")
	}
}

func (ctl *PermissionInternalController) parse(ctx *gin.Context) (
	user primitive.Account, resource domain.Resource, err error,
) {
	req := reqToCheckPermission{}
	if err = ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	user, index, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if resource, err = ctl.repo.GetByName(&index); err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			commonctl.SendError(
				ctx,
				allerror.NewNotFound(allerror.ErrorCodeRepoNotFound, "no repo"),
			)
		} else {
			commonctl.SendError(ctx, err)
		}
	}

	return
}
