/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	activityapp "github.com/openmerlin/merlin-server/activity/app"
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
	sl middleware.SecurityLog,
	rl middleware.RateLimiter,
) {
	m := ctl.userMiddleWare

	r.POST(`/v1/space`, m.Write,
		userctl.CheckMail(ctl.userMiddleWare, ctl.user, sl), l.Write, rl.CheckLimit, ctl.Create)
	r.DELETE("/v1/space/:id", m.Write,
		userctl.CheckMail(ctl.userMiddleWare, ctl.user, sl), l.Write, rl.CheckLimit, ctl.Delete)
	r.PUT("/v1/space/:id", m.Write,
		userctl.CheckMail(ctl.userMiddleWare, ctl.user, sl), l.Write, rl.CheckLimit, ctl.Update)
}

// SpaceController is a struct that contains the necessary dependencies for handling space-related operations.
type SpaceController struct {
	appService          app.SpaceAppService
	variableService     app.SpaceVariableService
	secretService       app.SpaceSecretService
	userMiddleWare      middleware.UserMiddleWare
	user                userapp.UserService
	rateLimitMiddleWare middleware.RateLimiter
	activity            activityapp.ActivityAppService
}

// @Summary  Create
// @Description  create space
// @Tags     Space
// @Param    body  body      reqToCreateSpace  true  "body of creating space"
// @Accept   json
// @Security Bearer
// @Success  201   {object}  commonctl.ResponseData{data=string,msg=string,code=string}
// @Router   /v1/space [post]
func (ctl *SpaceController) Create(ctx *gin.Context) {
	middleware.SetAction(ctx, "create space")

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
// @Param    id    path  string        true  "id of space" MaxLength(20)
// @Accept   json
// @Security Bearer
// @Success  204
// @Router   /v1/space/{id} [delete]
func (ctl *SpaceController) Delete(ctx *gin.Context) {
	middleware.SetAction(ctx, fmt.Sprintf("delete space of %s", ctx.Param("id")))

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
// @Param    id    path  string            true  "id of space" MaxLength(20)
// @Param    body  body  reqToUpdateSpace  true  "body of updating space"
// @Accept   json
// @Security Bearer
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Router   /v1/space/{id} [put]
func (ctl *SpaceController) Update(ctx *gin.Context) {
	middleware.SetAction(ctx, fmt.Sprintf("update space of %s", ctx.Param("id")))

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

// @Summary  Disable space
// @Description  disable space
// @Tags     Space
// @Param    id    path  string            true  "id of space" MaxLength(20)
// @Param    body  body  reqToDisableSpace  true  "body of disable space"
// @Accept   json
// @Security Bearer
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Router   /v1/space/{id}/disable [put]
func (ctl *SpaceController) Disable(ctx *gin.Context) {
	middleware.SetAction(ctx, fmt.Sprintf("disable space of %s", ctx.Param("id")))

	req := reqToDisableSpace{}
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

	action, err := ctl.appService.Disable(
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
