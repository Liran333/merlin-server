package server

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/openmerlin/merlin-server/api"
	"github.com/openmerlin/merlin-server/common/controller/middleware/operationlog"
	"github.com/openmerlin/merlin-server/config"
	sessionctl "github.com/openmerlin/merlin-server/session/controller"
)

func setRouterOfWeb(prefix string, engine *gin.Engine, cfg *config.Config, services *allServices) {
	api.SwaggerInfo.BasePath = prefix

	rg := engine.Group(api.SwaggerInfo.BasePath)

	services.userMiddleWare = sessionctl.WebAPIMiddleware(services.sessionApp)
	services.operationLog = operationlog.OperationLog(services.userMiddleWare)

	// set routers
	setRouterOfOrg(rg, cfg, services)

	setRouterOfUser(rg, cfg, services)

	setRouterOfSession(rg, services)

	setRouterOfModelWeb(rg, services)

	setRouterOfSpaceWeb(rg, services)

	setRouterOfSpaceAppWeb(rg, services)

	setRouterOfCodeRepoFile(rg, services)

	rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
