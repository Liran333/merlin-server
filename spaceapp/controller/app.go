/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides the controllers for handling HTTP requests and managing the application's business logic.
package controller

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	spacedomain "github.com/openmerlin/merlin-server/space/domain"
	"github.com/openmerlin/merlin-server/spaceapp/app"
	"github.com/openmerlin/merlin-server/spaceapp/domain"
)

// AddRouterForSpaceappWebController adds a router for the SpaceAppWebController to the given gin.RouterGroup.
func AddRouterForSpaceappWebController(
	r *gin.RouterGroup,
	s app.SpaceappAppService,
	m middleware.UserMiddleWare,
	l middleware.RateLimiter,
) {
	ctl := SpaceAppWebController{
		appService:          s,
		userMiddleWare:      m,
		rateLimitMiddleWare: l,
	}

	r.GET("/v1/space-app/:owner/:name", m.Optional, l.CheckLimit, ctl.Get)
	r.GET("/v1/space-app/:owner/:name/buildlog/realtime", m.Read, l.CheckLimit, ctl.GetRealTimeBuildLog)
	r.GET("/v1/space-app/:owner/:name/spacelog/realtime", m.Read, l.CheckLimit, ctl.GetRealTimeSpaceLog)
	r.POST("/v1/space-app/:owner/:name/restart", m.Optional, l.CheckLimit, ctl.Restart)
}

// SpaceAppWebController is a struct that represents the web controller for the space app.
type SpaceAppWebController struct {
	appService          app.SpaceappAppService
	userMiddleWare      middleware.UserMiddleWare
	rateLimitMiddleWare middleware.RateLimiter
}

// @Summary  Get
// @Description  get space app
// @Tags     SpaceAppWeb
// @Param    owner  path  string  true  "owner of space"
// @Param    name   path  string  true  "name of space"
// @Accept   json
// @Success  200  {object}  app.SpaceAppDTO
// @Router   /v1/space-app/{owner}/{name} [get]
func (ctl *SpaceAppWebController) Get(ctx *gin.Context) {
	index, err := ctl.parseIndex(ctx)
	if err != nil {
		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	if dto, err := ctl.appService.GetByName(user, &index); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, &dto)
	}
}

func (ctl *SpaceAppWebController) parseIndex(ctx *gin.Context) (index spacedomain.SpaceIndex, err error) {
	index.Owner, err = primitive.NewAccount(ctx.Param("owner"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	index.Name, err = primitive.NewMSDName(ctx.Param("name"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
	}

	return
}

// @Summary  GetBuildLog
// @Description  get space app real-time build log
// @Tags     SpaceAppWeb
// @Param    owner  path  string  true  "owner of space"
// @Param    name   path  string  true  "name of space"
// @Accept   json
// @Success  200  {object}  app.SpaceAppDTO
// @Router   /v1/space-app/{owner}/{name}/buildlog/realtime [get]
func (ctl *SpaceAppWebController) GetRealTimeBuildLog(ctx *gin.Context) {
	index, err := ctl.parseIndex(ctx)
	if err != nil {
		return
	}
	user := ctl.userMiddleWare.GetUser(ctx)

	spaceApp, err := ctl.appService.GetByName(user, &index)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}
	if spaceApp.BuildLogURL == "" {
		commonctl.SendBadRequestParam(ctx, fmt.Errorf("space app is not building"))
		return
	}

	streamWrite := func(doOnce func() ([]byte, error)) {
		ctx.Stream(func(w io.Writer) bool {
			done, err := doOnce()
			if err != nil {
				return false
			}
			ctx.SSEvent("message", string(done))
			return true
		})
	}

	params := domain.StreamParameter{
		Token:     config.SSEToken,
		StreamUrl: spaceApp.BuildLogURL,
	}
	cmd := &domain.SeverSentStream{
		Parameter:   params,
		Ctx:         ctx,
		StreamWrite: streamWrite,
	}

	if err := ctl.appService.GetRequestDataStream(cmd); err != nil {
		commonctl.SendError(ctx, err)
	}

}

// @Summary  GetSpaceLog
// @Description  get space app real-time space log
// @Tags     SpaceAppWeb
// @Param    owner  path  string  true  "owner of space"
// @Param    name   path  string  true  "name of space"
// @Accept   json
// @Success  200  {object}  app.SpaceAppDTO
// @Router   /v1/space-app/:owner/:name/spacelog/realtime [get]
func (ctl *SpaceAppWebController) GetRealTimeSpaceLog(ctx *gin.Context) {
	index, err := ctl.parseIndex(ctx)
	if err != nil {
		return
	}
	user := ctl.userMiddleWare.GetUser(ctx)

	spaceApp, err := ctl.appService.GetByName(user, &index)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}
	if spaceApp.AppLogURL == "" {
		commonctl.SendBadRequestParam(ctx, fmt.Errorf("space app is not serving"))
		return
	}

	streamWrite := func(doOnce func() ([]byte, error)) {
		ctx.Stream(func(w io.Writer) bool {
			done, err := doOnce()
			if err != nil {
				return false
			}
			ctx.SSEvent("message", string(done))
			return true
		})
	}

	params := domain.StreamParameter{
		Token:     config.SSEToken,
		StreamUrl: spaceApp.AppLogURL,
	}
	cmd := &domain.SeverSentStream{
		Parameter:   params,
		Ctx:         ctx,
		StreamWrite: streamWrite,
	}

	if err := ctl.appService.GetRequestDataStream(cmd); err != nil {
		commonctl.SendError(ctx, err)
	}

}

// @Summary  Post
// @Description  restart space app
// @Tags     SpaceAppWeb
// @Param    owner  path  string  true  "owner of space"
// @Param    name   path  string  true  "name of space"
// @Accept   json
// @Success  201   {object}  commonctl.ResponseData
// @Router   /v1/space-app/{owner}/{name}/restart [post]
func (ctl *SpaceAppWebController) Restart(ctx *gin.Context) {
	index, err := ctl.parseIndex(ctx)
	if err != nil {
		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	if err := ctl.appService.RestartSpaceApp(user, &index); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, "successfully")
	}
}
