package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

var restfulAPIInstance *restfulAPI

func RestfulAPI() *restfulAPI {
	return restfulAPIInstance
}

func Init() {
	// init both restfulAPI and webAPI
}

// restfulAPI
type restfulAPI struct {
}

func (m *restfulAPI) Write(ctx *gin.Context) {
}

func (m *restfulAPI) Optional(ctx *gin.Context) {
}

func (m *restfulAPI) GetUser(ctx *gin.Context) primitive.Account {
	return nil
}
