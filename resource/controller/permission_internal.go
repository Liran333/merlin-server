package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/coderepo/domain"
	commonapp "github.com/openmerlin/merlin-server/common/app"
	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	modelrepo "github.com/openmerlin/merlin-server/models/domain/repository"
	spacerepo "github.com/openmerlin/merlin-server/space/domain/repository"
)

var errorNoPermission = allerror.NewNoPermission("no permission")

func AddRouteForResourcePermissionInternalController(
	r *gin.RouterGroup,
	s commonapp.ResourcePermissionAppService,
	model modelrepo.ModelRepositoryAdapter,
	space spacerepo.SpaceRepositoryAdapter,
	m middleware.UserMiddleWare,
) {

	ctl := PermissionInternalController{
		ps:    s,
		model: model,
		space: space,
	}

	r.POST(`/v1/resource/permission/update`, m.Write, ctl.Update)
	r.POST(`/v1/resource/permission/delete`, m.Write, ctl.Delete)
	r.POST(`/v1/resource/permission/read`, m.Write, ctl.Read)
}

type PermissionInternalController struct {
	ps    commonapp.ResourcePermissionAppService
	model modelrepo.ModelRepositoryAdapter
	space spacerepo.SpaceRepositoryAdapter
}

// @Summary  Update
// @Description  check if can update resource
// @Tags     Permission
// @Param    body  body  reqToCheckPermission  true  "body of request"
// @Accept   json
// @Success  201   {object}  commonctl.ResponseData
// @Security Internal
// @Router   /v1/resource/permission/update [post]
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

// @Summary  Delete
// @Description  check if can delete resource
// @Tags     Permission
// @Param    body  body  reqToCheckPermission  true  "body of request"
// @Accept   json
// @Success  201   {object}  commonctl.ResponseData
// @Security Internal
// @Router   /v1/resource/permission/delete [post]
func (ctl *PermissionInternalController) Delete(ctx *gin.Context) {
	user, r, err := ctl.parse(ctx)
	if err != nil {
		return
	}

	if err := ctl.ps.CanDelete(user, r); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, "successfully")
	}
}

// @Summary  Read
// @Description  check if can read resource
// @Tags     Permission
// @Param    body  body  reqToCheckPermission  true  "body of request"
// @Accept   json
// @Success  201   {object}  commonctl.ResponseData
// @Security Internal
// @Router   /v1/resource/permission/read [post]
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
	user primitive.Account, resource commonapp.Resource, err error,
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

	if resource, err = ctl.getResource(index); err != nil {
		commonctl.SendError(ctx, err)
	}

	return
}

func (ctl *PermissionInternalController) getResource(index domain.CodeRepoIndex) (commonapp.Resource, error) {
	if r, err := ctl.model.FindByName(&index); err == nil {
		return &r, nil
	}

	if r, err := ctl.space.FindByName(&index); err == nil {
		return &r, nil
	}

	return nil, errorNoPermission
}
