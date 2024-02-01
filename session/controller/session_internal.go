package controller

import (
	"github.com/gin-gonic/gin"

	sdk "github.com/openmerlin/merlin-sdk/session"
	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/session/app"
)

func AddRouterForSessionInternalController(
	rg *gin.RouterGroup,
	s app.SessionAppService,
	m middleware.UserMiddleWare,
) {
	pc := SessionInternalController{
		s: s,
	}

	rg.PUT("/v1/session/check", m.Write, pc.CheckAndRefresh)
}

type SessionInternalController struct {
	s app.SessionAppService
}

// @Summary  CheckAndRefresh
// @Description  check and refresh session
// @Tags     Session
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

func cmdToCheck(req *sdk.RequestToCheckAndRefresh) (cmd app.CmdToCheck, err error) {
	if cmd.LoginId, err = primitive.NewUUID(req.LoginId); err != nil {
		return
	}

	if cmd.CSRFToken, err = primitive.NewUUID(req.CSRFToken); err != nil {
		return
	}

	cmd.IP = req.IP
	cmd.UserAgent = primitive.CreateUserAgent(req.UserAgent)

	return
}
