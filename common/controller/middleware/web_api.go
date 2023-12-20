package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

var webAPIInstance *webAPI

func WebAPI() *webAPI {
	return webAPIInstance
}

type webAPI struct {
}

func (m *webAPI) Write(ctx *gin.Context) {
}

func (m *webAPI) Optional(ctx *gin.Context) {
}

func (m *webAPI) GetUser(ctx *gin.Context) primitive.Account {
	return nil
}
