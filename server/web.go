package server

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/openmerlin/merlin-server/api"
	"github.com/openmerlin/merlin-server/config"

	sessionctl "github.com/openmerlin/merlin-server/session/controller"
)

func setRouterOfWeb(prefix string, engine *gin.Engine, cfg *config.Config, services *allServices) {
	api.SwaggerInfo.BasePath = prefix

	rg := engine.Group(api.SwaggerInfo.BasePath)

	services.userMiddleWare = sessionctl.WebAPIMiddleware(services.sessionApp)

	// set routers
	setRouterOfOrg(rg, cfg, services)

	setRouterOfUser(rg, cfg, services)

	setRouterOfSession(rg, services)

	setRouterOfModelWeb(rg, services)

	setRouterOfSpaceWeb(rg, services)

	setRouterOfCodeRepoFile(rg, services)

	rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
