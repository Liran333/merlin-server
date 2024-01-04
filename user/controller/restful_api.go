package controller

import (
	"strings"

	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	userapp "github.com/openmerlin/merlin-server/user/app"
)

const (
	authHeader      = "Authorization"
	authTokenPrefix = "token "
	authTokenKey    = "user_id"
)

func RestfulAPI(app userapp.UserService) *restfulAPI {
	return &restfulAPI{app}
}

// restfulAPI
type restfulAPI struct {
	s userapp.UserService
}

func (m *restfulAPI) Write(ctx *gin.Context) {
	t, err := m.s.VerifyToken(getToken(ctx), primitive.NewWritePerm())
	if err != nil {
		commonctl.SendError(ctx, err)
		ctx.Abort()
	} else {
		ctx.Set(authTokenKey, t)
		ctx.Next()
	}
}

func (m *restfulAPI) Optional(ctx *gin.Context) {
	token := getToken(ctx)
	if token == "" {
		ctx.Next()
		return
	}

	t, err := m.s.VerifyToken(token, primitive.NewReadPerm())
	if err != nil {
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

func getToken(ctx *gin.Context) string {
	auth := ctx.GetHeader(authHeader)
	if auth == "" || !strings.HasPrefix(auth, authTokenPrefix) {
		return ""
	}

	return strings.TrimPrefix(auth, authTokenPrefix)
}
