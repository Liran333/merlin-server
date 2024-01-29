package server

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/openmerlin/merlin-server/api"
	"github.com/openmerlin/merlin-server/config"

	userctl "github.com/openmerlin/merlin-server/user/controller"
)

func setRouterOfRestful(prefix string, engine *gin.Engine, cfg *config.Config, services *allServices) {
	api.SwaggerInfo.BasePath = prefix

	rg := engine.Group(api.SwaggerInfo.BasePath)

	services.userMiddleWare = userctl.RestfulAPI(services.userApp)

	// set routers
	setRouterOfOrg(rg, cfg, services)

	setRouterOfUser(rg, cfg, services)

	setRouterOfCodeRepoFile(rg, services)

	setRouterOfModelRestful(rg, services)

	setRouterOfSpaceRestful(rg, services)

	setRouterOfBranchRestful(rg, services)

	rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
