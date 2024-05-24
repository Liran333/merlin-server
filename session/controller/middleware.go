/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/session/app"
)

const (
	userIdParsed    = "user_id"
	csrfTokenHeader = "Csrf-Token" // #nosec G101
)

var noUserError = errors.New("no user")

// WebAPIMiddleware creates a new instance of webAPIMiddleware with the given session and securityLog.
func WebAPIMiddleware(session app.SessionAppService, securityLog middleware.SecurityLog,
	cfg *Config) *webAPIMiddleware {
	return &webAPIMiddleware{
		session:     session,
		securityLog: securityLog,
		cfg:         cfg,
	}
}

type webAPIMiddleware struct {
	session     app.SessionAppService
	securityLog middleware.SecurityLog
	cfg         *Config
}

// Write is a middleware function that checks if the CSRF token is present in the request header.
func (m *webAPIMiddleware) Write(ctx *gin.Context) {
	m.must(ctx)
}

// Read is a middleware function that checks if the CSRF token is present in the request header.
func (m *webAPIMiddleware) Read(ctx *gin.Context) {
	m.must(ctx)
}

// CheckSession is a middleware function that checks if the session is present in the request header.
func (m *webAPIMiddleware) CheckSession(ctx *gin.Context) {
	m.mustSession(ctx)
}

// Optional is a middleware function that checks if the CSRF token is present in the request header.
func (m *webAPIMiddleware) Optional(ctx *gin.Context) {
	if v := ctx.GetHeader(csrfTokenHeader); v == "" {
		ctx.Next()
	} else {
		m.must(ctx)
	}
}

func (m *webAPIMiddleware) must(ctx *gin.Context) {
	if err := m.checkToken(ctx); err != nil {
		clearCookie(ctx, m.cfg.SessionDomain)
		commonctl.SendError(ctx, err)
		m.securityLog.Warn(ctx, err.Error())

		ctx.Abort()
	} else {
		ctx.Next()
	}
}

func (m *webAPIMiddleware) mustSession(ctx *gin.Context) {
	if err := m.checkSession(ctx); err != nil {
		clearCookie(ctx, m.cfg.SessionDomain)
		commonctl.SendError(ctx, err)
		m.securityLog.Warn(ctx, err.Error())

		ctx.Abort()
	} else {
		ctx.Next()
	}
}

// GetUser retrieves the user account from the context.
func (m *webAPIMiddleware) GetUser(ctx *gin.Context) primitive.Account {
	v, ok := ctx.Get(userIdParsed)
	if !ok {
		return nil
	}

	if r, ok := v.(primitive.Account); ok {
		return r
	}

	return nil
}

// GetUserAndExitIfFailed retrieves the user account from the context and exits if the user is not found.
func (m *webAPIMiddleware) GetUserAndExitIfFailed(ctx *gin.Context) primitive.Account {
	if v := m.GetUser(ctx); v != nil {
		return v
	}

	commonctl.SendError(ctx, noUserError)

	return nil
}

func (m *webAPIMiddleware) checkToken(ctx *gin.Context) error {
	csrfToken, err := m.parseCSRFToken(ctx)
	if err != nil {
		return err
	}

	sessionId, err := m.parseSessionId(ctx)
	if err != nil {
		return err
	}

	ip, err := commonctl.GetIp(ctx)
	if err != nil {
		return err
	}

	userAgent, err := commonctl.GetUserAgent(ctx)
	if err != nil {
		return err
	}

	user, newCSRF, err := m.session.CheckAndRefresh(&app.CmdToCheck{
		SessionDTO: app.SessionDTO{
			SessionId: sessionId,
			CSRFToken: csrfToken,
		},
		IP:        ip,
		UserAgent: userAgent,
	})
	if err != nil {
		return err
	}

	expiry := config.csrfTokenCookieExpiry()
	setCookieOfCSRFToken(ctx, newCSRF, m.cfg.SessionDomain, &expiry)

	ctx.Set(userIdParsed, user)

	return nil
}

func (m *webAPIMiddleware) checkSession(ctx *gin.Context) error {
	sessionId, err := m.parseSessionId(ctx)
	if err != nil {
		return err
	}

	user, err := m.session.CheckSession(&app.CmdToCheck{
		SessionDTO: app.SessionDTO{
			SessionId: sessionId,
		},
	})
	if err != nil {
		return err
	}

	ctx.Set(userIdParsed, user)

	return nil
}

func (m *webAPIMiddleware) parseCSRFToken(ctx *gin.Context) (primitive.RandomId, error) {
	v := ctx.GetHeader(csrfTokenHeader)
	if v == "" {
		return nil, allerror.New(
			allerror.ErrorCodeCSRFTokenMissing, "no csrf token", fmt.Errorf("no csrf token found"),
		)
	}

	id, err := primitive.ToRandomId(v)
	if err != nil {
		err = allerror.New(
			allerror.ErrorCodeCSRFTokenInvalid, "invalid csrf token", fmt.Errorf("invalid csrf token"),
		)
	}

	return id, err
}

func (m *webAPIMiddleware) parseSessionId(ctx *gin.Context) (primitive.RandomId, error) {
	v, err := commonctl.GetCookie(ctx, config.CookieSessionId)
	if err != nil {
		return nil, allerror.New(allerror.ErrorCodeCSRFTokenMissing, "no session id found", err)
	}

	sessionId, err := primitive.ToRandomId(v)
	if err != nil {
		err = allerror.New(allerror.ErrorCodeSessionIdInvalid, "not session id", err)
	}

	return sessionId, err
}

func (m *webAPIMiddleware) ClearCookieAfterRevokePrivacy(ctx *gin.Context) {
	clearCookie(ctx, m.cfg.SessionDomain)
}

func clearCookie(ctx *gin.Context, domain string) {
	expiry := time.Now().AddDate(0, 0, -1)
	setCookieOfSessionId(ctx, "", domain, &expiry)
	setCookieOfCSRFToken(ctx, "", domain, &expiry)
}
