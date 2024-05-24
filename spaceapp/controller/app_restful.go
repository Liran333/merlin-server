/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package controller provides the controllers for handling HTTP requests and managing the application's business logic.
package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/spaceapp/app"
)

// AddRouterForSpaceappRestfulController adds a router for the SpaceAppWebController to the given gin.RouterGroup.
func AddRouterForSpaceappRestfulController(
	r *gin.RouterGroup,
	s app.SpaceappAppService,
	m middleware.UserMiddleWare,
	t middleware.TokenMiddleWare,
	l middleware.RateLimiter,
) {
	ctl := SpaceRestfulController{
		SpaceAppController: SpaceAppController{
			appService:          s,
			userMiddleWare:      m,
			tokenMiddleWare:     t,
			rateLimitMiddleWare: l,
		},
	}

	addRouterForSpaceappController(r, &ctl.SpaceAppController, m, l)

}

// SpaceRestfulController is a struct that represents the restful controller for the space app.
type SpaceRestfulController struct {
	SpaceAppController
}
