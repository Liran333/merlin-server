package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

var instance *userCheckingMiddleware

func UserMiddleware() *userCheckingMiddleware {
	return instance
}

type userCheckingMiddleware struct {
}

func (m *userCheckingMiddleware) Must(ctx *gin.Context) {
	// TODO
	// 1. The token must be exist and valid, otherwise abor directly and
	// send allerror.ErrorCodeAccessTokenInvalid
	// 2. If token is valid, then parse the user bound to the token and save it
}

func (m *userCheckingMiddleware) Optional(ctx *gin.Context) {
	// TODO
	// 1. if token is not passed, ignore it.
	// 2. If token exists, it must be valid, otherwise abor directly and
	// send allerror.ErrorCodeAccessTokenInvalid
	// 3. If token is valid, then parse the user bound to the token and save it
}

func (m *userCheckingMiddleware) GetUser(ctx *gin.Context) primitive.Account {
	return nil
}
