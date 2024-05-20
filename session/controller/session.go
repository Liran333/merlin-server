/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	orgapp "github.com/openmerlin/merlin-server/organization/app"
	"github.com/openmerlin/merlin-server/session/app"
)

// AddRouterForSessionController adds routes for session controller to the given router group.
func AddRouterForSessionController(
	rg *gin.RouterGroup,
	s app.SessionAppService,
	l middleware.OperationLog,
	m middleware.UserMiddleWare,
	cfg *Config,
	d orgapp.PrivilegeOrg,
) {
	pc := SessionController{
		s:       s,
		cfg:     cfg,
		disable: d,
	}

	rg.POST("/v1/session", l.Write, pc.Login)
	rg.PUT("/v1/session", m.Write, l.Write, pc.Logout)
}

// SessionController is a struct that holds the session app service.
type SessionController struct {
	s       app.SessionAppService
	cfg     *Config
	disable orgapp.PrivilegeOrg
}

// @Summary  Login
// @Description  login
// @Tags     Session
// @Param    body  body  reqToLogin  true  "body of login"
// @Accept   json
// @Success  201   {object}  commonctl.ResponseData{data=app.UserDTO,msg=string,code=string}
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

	sessionExpiry := config.sessionCookieExpiry()
	setCookieOfSessionId(ctx, session.SessionId.RandomId(), ctl.cfg.SessionDomain, &sessionExpiry)

	expiry := config.csrfTokenCookieExpiry()
	setCookieOfCSRFToken(ctx, session.CSRFToken.RandomId(), ctl.cfg.SessionDomain, &expiry)

	u := app.UserInfoDTO{
		UserDTO: user,
	}
	if ctl.disable != nil {
		userAccount, _ := primitive.NewAccount(user.Name)
		err = ctl.disable.Contains(userAccount)
		if err != nil {
			u.IsDisableAdmin = false
		} else {
			u.IsDisableAdmin = true
		}
	}

	commonctl.SendRespOfGet(ctx, u)
}

// @Summary  Logout
// @Description  logout
// @Tags     Session
// @Accept   json
// @Success  202  {object}  commonctl.ResponseData{data=logoutInfo,msg=string,code=string}
// @Router   /v1/session [put]
func (ctl *SessionController) Logout(ctx *gin.Context) {
	v, err := commonctl.GetCookie(ctx, config.CookieSessionId)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	middleware.SetAction(ctx, "logout")

	sessionId, err := primitive.ToRandomId(v)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	idToken, err := ctl.s.Logout(sessionId)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	expiry := time.Now().AddDate(0, 0, -1)
	setCookieOfSessionId(ctx, "", ctl.cfg.SessionDomain, &expiry)
	setCookieOfCSRFToken(ctx, "", ctl.cfg.SessionDomain, &expiry)

	commonctl.SendRespOfPut(ctx, logoutInfo{IdToken: idToken})
}

func setCookieOfCSRFToken(ctx *gin.Context, value, domain string, expiry *time.Time) {
	commonctl.SetCookie(ctx, config.CookieCSRFToken, value, domain, false, expiry, http.SameSiteStrictMode)
	if config.LocalDomainCookie {
		commonctl.SetCookie(ctx, config.CookieCSRFToken, value, "", false, expiry, http.SameSiteLaxMode)
	}
}

func setCookieOfSessionId(ctx *gin.Context, value, domain string, expiry *time.Time) {
	commonctl.SetCookie(ctx, config.CookieSessionId, value, domain, true, expiry, http.SameSiteLaxMode)
	if config.LocalDomainCookie {
		commonctl.SetCookie(ctx, config.CookieSessionId, value, "", true, expiry, http.SameSiteLaxMode)
	}
}
