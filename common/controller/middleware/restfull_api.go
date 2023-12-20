package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

var restfullAPIInstance *restfullAPI

func RestfullAPI() *restfullAPI {
	return restfullAPIInstance
}

func Init() {
	// init both restfullAPI and webAPI
}

// restfullAPI
type restfullAPI struct {
}

func (m *restfullAPI) Write(ctx *gin.Context) {
}

func (m *restfullAPI) Optional(ctx *gin.Context) {
}

func (m *restfullAPI) GetUser(ctx *gin.Context) primitive.Account {
	return nil
}
