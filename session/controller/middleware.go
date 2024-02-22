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
	csrfTokenHeader = "CSRF-TOKEN" // #nosec G101
)

var noUserError = errors.New("no user")

func WebAPIMiddleware(session app.SessionAppService, securityLog middleware.SecurityLog) *webAPIMiddleware {
	return &webAPIMiddleware{
		session,
		securityLog,
	}
}

type webAPIMiddleware struct {
	session     app.SessionAppService
	securityLog middleware.SecurityLog
}

func (m *webAPIMiddleware) Write(ctx *gin.Context) {
	m.must(ctx)
}

func (m *webAPIMiddleware) Read(ctx *gin.Context) {
	m.must(ctx)
}

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

	loginId, err := m.parseLoginId(ctx)
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
			LoginId:   loginId,
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

func (m *webAPIMiddleware) parseCSRFToken(ctx *gin.Context) (primitive.UUID, error) {
	v := ctx.GetHeader(csrfTokenHeader)
	if v == "" {
		return primitive.UUID{}, allerror.New(
			allerror.ErrorCodeCSRFTokenMissing, "no csrf token",
		)
	}

	csrfToken, err := primitive.NewUUID(v)
	if err != nil {
		err = allerror.New(allerror.ErrorCodeCSRFTokenInvalid, "not uuid")
	}

	return csrfToken, err
}

func (m *webAPIMiddleware) parseLoginId(ctx *gin.Context) (primitive.UUID, error) {
	v, err := commonctl.GetCookie(ctx, cookieLoginId)
	if err != nil {
		return primitive.UUID{}, err
	}

	loginId, err := primitive.NewUUID(v)
	if err != nil {
		err = allerror.New(allerror.ErrorCodeLoginIdInvalid, "not uuid")
	}

	return loginId, err
}
