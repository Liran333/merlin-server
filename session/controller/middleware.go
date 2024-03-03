/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"errors"

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
func WebAPIMiddleware(session app.SessionAppService, securityLog middleware.SecurityLog) *webAPIMiddleware {
	return &webAPIMiddleware{
		session:     session,
		securityLog: securityLog,
	}
}

type webAPIMiddleware struct {
	session     app.SessionAppService
	securityLog middleware.SecurityLog
}

// Write is a middleware function that checks if the CSRF token is present in the request header.
func (m *webAPIMiddleware) Write(ctx *gin.Context) {
	m.must(ctx)
}

// Read is a middleware function that checks if the CSRF token is present in the request header.
func (m *webAPIMiddleware) Read(ctx *gin.Context) {
	m.must(ctx)
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
	setCookieOfCSRFToken(ctx, newCSRF, &expiry)

	ctx.Set(userIdParsed, user)

	return nil
}

func (m *webAPIMiddleware) parseCSRFToken(ctx *gin.Context) (primitive.RandomId, error) {
	v := ctx.GetHeader(csrfTokenHeader)
	if v == "" {
		return nil, allerror.New(
			allerror.ErrorCodeCSRFTokenMissing, "no csrf token",
		)
	}

	return primitive.ToRandomId(v)
}

func (m *webAPIMiddleware) parseSessionId(ctx *gin.Context) (primitive.RandomId, error) {
	v, err := commonctl.GetCookie(ctx, cookieSessionId)
	if err != nil {
		return nil, err
	}

	sessionId, err := primitive.ToRandomId(v)
	if err != nil {
		err = allerror.New(allerror.ErrorCodeSessionIdInvalid, "not session id")
	}

	return sessionId, err
}
