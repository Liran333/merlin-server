/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/space/app"
	"github.com/openmerlin/merlin-server/space/domain"
	userapp "github.com/openmerlin/merlin-server/user/app"
	userctl "github.com/openmerlin/merlin-server/user/controller"
)

func addRouteForSpaceController(
	r *gin.RouterGroup,
	ctl *SpaceController,
	l middleware.OperationLog,
) {
	m := ctl.userMiddleWare

	r.POST(`/v1/space`, m.Write, userctl.CheckMail(ctl.userMiddleWare, ctl.user), l.Write, ctl.Create)
	r.DELETE("/v1/space/:id", m.Write, userctl.CheckMail(ctl.userMiddleWare, ctl.user), l.Write, ctl.Delete)
	r.PUT("/v1/space/:id", m.Write, userctl.CheckMail(ctl.userMiddleWare, ctl.user), l.Write, ctl.Update)
}

// SpaceController is a struct that contains the necessary dependencies for handling space-related operations.
type SpaceController struct {
	appService     app.SpaceAppService
	userMiddleWare middleware.UserMiddleWare
	user           userapp.UserService
}

// @Summary  Create
// @Description  create space
// @Tags     Space
// @Param    body  body      reqToCreateSpace  true  "body of creating space"
// @Accept   json
// @Security Bearer
// @Success  201   {object}  commonctl.ResponseData
// @Router   /v1/space [post]
func (ctl *SpaceController) Create(ctx *gin.Context) {
	req := reqToCreateSpace{}
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

	if v, err := ctl.appService.Create(user, &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, v)
	}
}

// @Summary  Delete
// @Description  delete space
// @Tags     Space
// @Param    id    path  string        true  "id of space"
// @Accept   json
// @Security Bearer
// @Success  204
// @Router   /v1/space/{id} [delete]
func (ctl *SpaceController) Delete(ctx *gin.Context) {
	user := ctl.userMiddleWare.GetUser(ctx)

	spaceId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	action, err := ctl.appService.Delete(user, spaceId)

	middleware.SetAction(ctx, action)

	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfDelete(ctx)
	}
}

// @Summary  Update
// @Description  update space
// @Tags     Space
// @Param    id    path  string            true  "id of space"
// @Param    body  body  reqToUpdateSpace  true  "body of updating space"
// @Accept   json
// @Security Bearer
// @Success  202   {object}  commonctl.ResponseData
// @Router   /v1/space/{id} [put]
func (ctl *SpaceController) Update(ctx *gin.Context) {
	req := reqToUpdateSpace{}
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

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

	action, err := ctl.appService.Update(
		ctl.userMiddleWare.GetUser(ctx),
		spaceId, &cmd,
	)

	middleware.SetAction(ctx, fmt.Sprintf("%s, set %s", action, req.action()))

	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

func (ctl *SpaceController) parseIndex(ctx *gin.Context) (index domain.SpaceIndex, err error) {
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
