package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/session/app"
)

const (
	cookieLoginId   = "login_id"
	cookieCSRFToken = "csrf_token"
)

func AddRouterForSessionController(
	rg *gin.RouterGroup,
	s app.SessionAppService,
	l middleware.OperationLog,
	m middleware.UserMiddleWare,
) {
	pc := SessionController{
		s: s,
	}

	rg.POST("/v1/session", l.Write, pc.Login)
	rg.PUT("/v1/session", m.Write, l.Write, pc.Logout)
}

type SessionController struct {
	s app.SessionAppService
}

// @Summary  Login
// @Description  login
// @Tags     Session
// @Param    body  body  reqToLogin  true  "body of login"
// @Accept   json
// @Success  201   {object}  app.UserDTO
// @Router   /v1/session [post]
func (ctl *SessionController) Login(ctx *gin.Context) {
	var req reqToLogin

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	middleware.SetAction(ctx, "login")

	cmd, err := req.toCmd(ctx)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	session, user, err := ctl.s.Login(&cmd)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	middleware.SetAction(ctx, fmt.Sprintf("%s login", user.Name))

	setCookieOfLoginId(ctx, session.LoginId.String(), nil)

	expiry := config.csrfTokenCookieExpiry()
	setCookieOfCSRFToken(ctx, session.CSRFToken.String(), &expiry)

	commonctl.SendRespOfGet(ctx, user)
}

// @Summary  Logout
// @Description  logout
// @Tags     Session
// @Accept   json
// @Success  202  {object}  logoutInfo
// @Router   /v1/session [put]
func (ctl *SessionController) Logout(ctx *gin.Context) {
	v, err := commonctl.GetCookie(ctx, cookieLoginId)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	middleware.SetAction(ctx, "logout")

	loginId, err := primitive.NewUUID(v)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	idToken, err := ctl.s.Logout(loginId)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	expiry := time.Now().AddDate(0, 0, -1)
	setCookieOfLoginId(ctx, "", &expiry)
	setCookieOfCSRFToken(ctx, "", &expiry)

	commonctl.SendRespOfPut(ctx, logoutInfo{idToken})
}

func setCookieOfLoginId(ctx *gin.Context, value string, expiry *time.Time) {
	commonctl.SetCookie(ctx, cookieLoginId, value, true, expiry)
}

func setCookieOfCSRFToken(ctx *gin.Context, value string, expiry *time.Time) {
	commonctl.SetCookie(ctx, cookieCSRFToken, value, false, expiry)
}
