package server

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/openmerlin/merlin-server/api"
	"github.com/openmerlin/merlin-server/config"

	"github.com/openmerlin/merlin-server/common/controller/middleware/internalservice"
)

func setRouterOfInternal(prefix string, engine *gin.Engine, cfg *config.Config, services *allServices) {
	api.SwaggerInfo.BasePath = prefix

	rg := engine.Group(api.SwaggerInfo.BasePath)

	services.userMiddleWare = internalservice.NewAPIMiddleware()

	// set routers
	setRouterOfSessionInternal(rg, services)
	setInternalRouterOfUser(rg, cfg, services)

	setRouterOfSpaceInternal(rg, services)

	setRouterOfModelInternal(rg, services)

	setRouterOfSpaceAppInternal(rg, services, cfg)

	rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
