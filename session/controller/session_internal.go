/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"github.com/gin-gonic/gin"
	sdk "github.com/openmerlin/merlin-sdk/session"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/session/app"
)

// AddRouterForSessionInternalController adds routes for session internal controller to the given router group.
func AddRouterForSessionInternalController(
	rg *gin.RouterGroup,
	s app.SessionAppService,
	l middleware.OperationLog,
	m middleware.UserMiddleWare,
) {
	pc := SessionInternalController{
		s: s,
	}

	rg.PUT("/v1/session/check", m.Write, l.Write, pc.CheckAndRefresh)
	rg.DELETE("/v1/session/clear", m.Write, l.Write, pc.Clear)
}

// SessionInternalController is a struct that holds the session app service.
type SessionInternalController struct {
	s app.SessionAppService
}

// @Summary  CheckAndRefresh
// @Description  check and refresh session
// @Tags     SessionInternal
// @Accept   json
// @Success  202
// @Router   /v1/session/check [put]
func (ctl *SessionInternalController) CheckAndRefresh(ctx *gin.Context) {
	var req sdk.RequestToCheckAndRefresh

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := cmdToCheck(&req)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if user, token, err := ctl.s.CheckAndRefresh(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, sdk.ResponseToCheckAndRefresh{
			User:      user.Account(),
			CSRFToken: token,
		})
	}
}

// @Summary  Clear session by session id
// @Description  Clear session when it expired
// @Tags     SessionInternal
// @Accept   json
// @Success  204
// @Router   /v1/session/clear [delete]
func (ctl *SessionInternalController) Clear(ctx *gin.Context) {
	var req sdk.RequestToClear

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	sessionId, err := primitive.ToRandomId(req.SessionId)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if err = ctl.s.Clear(sessionId); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfDelete(ctx)
	}
}

func cmdToCheck(req *sdk.RequestToCheckAndRefresh) (cmd app.CmdToCheck, err error) {
	if cmd.SessionId, err = primitive.ToRandomId(req.SessionId); err != nil {
		return
	}

	if cmd.CSRFToken, err = primitive.ToRandomId(req.CSRFToken); err != nil {
		return
	}

	cmd.IP = req.IP
	cmd.UserAgent = primitive.CreateUserAgent(req.UserAgent)

	return
}
