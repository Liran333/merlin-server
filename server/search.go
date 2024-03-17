/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package server

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/search/app"
	"github.com/openmerlin/merlin-server/search/controller"

	"github.com/openmerlin/merlin-server/search/infrastructure/resourceadapterimpl"

	"github.com/openmerlin/merlin-server/models/infrastructure/modelrepositoryadapter"
	"github.com/openmerlin/merlin-server/space/infrastructure/spacerepositoryadapter"
)

func setRouterOfSearchWeb(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouteForSearchWebController(
		rg,
		app.NewSearchAppService(
			resourceadapterimpl.NewSearchRepositoryAdapter(
				modelrepositoryadapter.ModelAdapter(),
				spacerepositoryadapter.SpaceAdapter(),
				services.userRepo,
			),
		),
		services.operationLog,
		services.userMiddleWare,
	)
}
