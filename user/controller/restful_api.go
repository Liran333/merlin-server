package controller

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	userapp "github.com/openmerlin/merlin-server/user/app"
)

const (
	authHeader      = "Authorization"
	authTokenPrefix = "Bearer "
	authTokenKey    = "user_id"
)

var errNoUserError = errors.New("no user")

func RestfulAPI(app userapp.UserService) *restfulAPI {
	return &restfulAPI{app}
}

// restfulAPI
type restfulAPI struct {
	s userapp.UserService
}

func (m *restfulAPI) Write(ctx *gin.Context) {
	m.check(ctx, false, primitive.NewWritePerm())
}

func (m *restfulAPI) Read(ctx *gin.Context) {
	m.check(ctx, false, primitive.NewReadPerm())
}

func (m *restfulAPI) Optional(ctx *gin.Context) {
	m.check(ctx, true, primitive.NewReadPerm())
}

func (m *restfulAPI) check(ctx *gin.Context, ignore bool, permission primitive.TokenPerm) {
	token := getToken(ctx)
	if token == "" {
		if ignore {
			ctx.Next()
		} else {
			commonctl.SendError(ctx, allerror.New(
				allerror.ErrorCodeAccessTokenInvalid, "missing token",
			))

			ctx.Abort()
		}

		return
	}

	if t, err := m.s.VerifyToken(token, permission); err != nil {
		commonctl.SendError(ctx, err)

		ctx.Abort()
	} else {
		ctx.Set(authTokenKey, t.Account)

		ctx.Next()
	}
}

func (m *restfulAPI) GetUser(ctx *gin.Context) primitive.Account {
	u, ok := ctx.Get(authTokenKey)
	if !ok {
		return nil
	}

	t, ok := u.(string)
	if !ok {
		return nil
	}

	return primitive.CreateAccount(t)
}

func (m *restfulAPI) GetUserAndExitIfFailed(ctx *gin.Context) primitive.Account {
	if v := m.GetUser(ctx); v != nil {
		return v
	}

	commonctl.SendError(ctx, errNoUserError)

	return nil
}

func getToken(ctx *gin.Context) string {
	auth := ctx.GetHeader(authHeader)
	if auth == "" || !strings.HasPrefix(auth, authTokenPrefix) {
		return ""
	}

	return strings.TrimPrefix(auth, authTokenPrefix)
}
