package internalservice

import (
	"errors"

	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

const tokenHeader = "TOKEN" // #nosec G101

var noUserError = errors.New("no user")

func NewAPIMiddleware(securityLog middleware.SecurityLog) *internalServiceAPIMiddleware {
	return &internalServiceAPIMiddleware{
		securityLog: securityLog,
	}
}

// internalServiceAPIMiddleware
type internalServiceAPIMiddleware struct {
	securityLog middleware.SecurityLog
}

func (m *internalServiceAPIMiddleware) Write(ctx *gin.Context) {
	m.must(ctx)
}

func (m *internalServiceAPIMiddleware) Read(ctx *gin.Context) {
	m.must(ctx)
}

func (m *internalServiceAPIMiddleware) Optional(ctx *gin.Context) {
	if v := ctx.GetHeader(tokenHeader); v == "" {
		ctx.Next()
	} else {
		m.must(ctx)
	}
}

func (m *internalServiceAPIMiddleware) must(ctx *gin.Context) {
	if err := m.checkToken(ctx); err != nil {
		commonctl.SendError(ctx, err)
		m.securityLog.Warn(ctx, err.Error())

		ctx.Abort()
	} else {
		ctx.Next()
	}
}

func (m *internalServiceAPIMiddleware) GetUser(ctx *gin.Context) primitive.Account {
	return nil
}

func (m *internalServiceAPIMiddleware) GetUserAndExitIfFailed(ctx *gin.Context) primitive.Account {
	commonctl.SendError(ctx, noUserError)

	return nil
}

func (m *internalServiceAPIMiddleware) checkToken(ctx *gin.Context) error {
	rawToken := ctx.GetHeader(tokenHeader)
	calcTokenHash, err := commonctl.EncodeToken(rawToken, config.Salt)
	if err != nil {
		return allerror.New(
			allerror.ErrorCodeAccessTokenInvalid, "check token failed",
		)
	}

	if calcTokenHash != config.TokenHash {
		return allerror.New(
			allerror.ErrorCodeAccessTokenInvalid, "invalid token",
		)
	}

	return nil
}
