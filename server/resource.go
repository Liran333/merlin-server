package server

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/models/infrastructure/modelrepositoryadapter"
	"github.com/openmerlin/merlin-server/resource/controller"
	"github.com/openmerlin/merlin-server/space/infrastructure/spacerepositoryadapter"
)

func setRouterOfResourcePermissionInternal(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouteForResourcePermissionInternalController(
		rg,
		services.permissionApp,
		modelrepositoryadapter.ModelAdapter(),
		spacerepositoryadapter.SpaceAdapter(),
		services.userMiddleWare,
	)
}
