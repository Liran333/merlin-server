/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package controller provides the controllers for handling HTTP requests and managing the application's business logic.
package controller

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

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
	t middleware.TokenMiddleWare,
	l middleware.RateLimiter,
) {
	ctl := SpaceAppWebController{
		SpaceAppController: SpaceAppController{
			appService:          s,
			userMiddleWare:      m,
			tokenMiddleWare:     t,
			rateLimitMiddleWare: l,
		},
	}

	addRouterForSpaceappController(r, &ctl.SpaceAppController, m, l)

	r.GET("/v1/space-app/:owner/:name", m.Optional, l.CheckLimit, ctl.Get)
	r.GET("/v1/space-app/:owner/:name/buildlog/realtime", m.Read, l.CheckLimit, ctl.GetRealTimeBuildLog)
	r.GET("/v1/space-app/:owner/:name/spacelog/realtime", m.Read, l.CheckLimit, ctl.GetRealTimeSpaceLog)
	r.GET("/v1/space-app/:owner/:name/read", t.CheckSession, l.CheckLimit, ctl.CanRead)
	r.GET("/v1/space-app/:owner/:name/buildlog/complete", m.Read, l.CheckLimit, ctl.GetBuildLogs)
}

// SpaceAppWebController is a struct that represents the web controller for the space app.
type SpaceAppWebController struct {
	SpaceAppController
}

// @Summary  Get
// @Description  get space app
// @Tags     SpaceAppWeb
// @Param    owner  path  string  true  "owner of space" MaxLength(40)
// @Param    name   path  string  true  "name of space" MaxLength(100)
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=app.SpaceAppDTO,msg=string,code=string}
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

// @Summary  GetBuildLogs
// @Description  get space app complete buid logs
// @Tags     SpaceAppWeb
// @Param    id  path  string  true  "space app id"
// @Accept   json
// @Success  200  {object}  app.BuildLogsDTO
// @Router   /v1/space-app/{owner}/{name}/buildlog/complete [get]
func (ctl *SpaceAppWebController) GetBuildLogs(ctx *gin.Context) {
	index, err := ctl.parseIndex(ctx)
	if err != nil {
		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	if dto, err := ctl.appService.GetBuildLogs(user, &index); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, &dto)
	}
}

// parseIndex parses the index from the request.
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
// @Param    owner  path  string  true  "owner of space" MaxLength(40)
// @Param    name   path  string  true  "name of space" MaxLength(100)
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=app.SpaceAppDTO,msg=string,code=string}
// @Router   /v1/space-app/{owner}/{name}/buildlog/realtime [get]
func (ctl *SpaceAppWebController) GetRealTimeBuildLog(ctx *gin.Context) {
	index, err := ctl.parseIndex(ctx)
	if err != nil {
		ctx.SSEvent("error", err.Error())
		return
	}
	user := ctl.userMiddleWare.GetUser(ctx)

	buildLog, err := ctl.appService.GetBuildLog(user, &index)
	if err != nil {
		logrus.Errorf("get build log err:%s", err)
		ctx.SSEvent("error", "get build log failed")
		return
	}

	streamWrite := func(doOnce func() ([]byte, error)) {
		ctx.Stream(func(w io.Writer) bool {
			done, err := doOnce()
			if err != nil {
				if err.Error() == "finish" {
					ctx.SSEvent("message", "")
				} else {
					logrus.Errorf("request build log err:%s", err)
					ctx.SSEvent("error", "request build log failed")
				}
				return false
			}
			if done != nil {
				ctx.SSEvent("message", string(done))
			}
			return true
		})
	}

	params := domain.StreamParameter{
		Token:     config.SSEToken,
		StreamUrl: buildLog,
	}
	cmd := &domain.SeverSentStream{
		Parameter:   params,
		Ctx:         ctx,
		StreamWrite: streamWrite,
	}

	if err := ctl.appService.GetRequestDataStream(cmd); err != nil {
		ctx.SSEvent("error", err.Error())
	}

}

// @Summary  GetSpaceLog
// @Description  get space app real-time space log
// @Tags     SpaceAppWeb
// @Param    owner  path  string  true  "owner of space" MaxLength(40)
// @Param    name   path  string  true  "name of space" MaxLength(100)
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=app.SpaceAppDTO,msg=string,code=string}
// @Router   /v1/space-app/:owner/:name/spacelog/realtime [get]
func (ctl *SpaceAppWebController) GetRealTimeSpaceLog(ctx *gin.Context) {
	index, err := ctl.parseIndex(ctx)
	if err != nil {
		ctx.SSEvent("error", err.Error())
		return
	}
	user := ctl.userMiddleWare.GetUser(ctx)

	spaceLog, err := ctl.appService.GetSpaceLog(user, &index)
	if err != nil {
		logrus.Errorf("get space log err:%s", err)
		ctx.SSEvent("error", "get space log failed")
		return
	}

	streamWrite := func(doOnce func() ([]byte, error)) {
		ctx.Stream(func(w io.Writer) bool {
			done, err := doOnce()
			if err != nil {
				if err.Error() == "finish" {
					ctx.SSEvent("message", "")
				} else {
					logrus.Errorf("request space log err:%s", err)
					ctx.SSEvent("error", "request space log failed")
				}
				return false
			}
			if done != nil {
				ctx.SSEvent("message", string(done))
			}
			return true
		})
	}

	params := domain.StreamParameter{
		Token:     config.SSEToken,
		StreamUrl: spaceLog,
	}
	cmd := &domain.SeverSentStream{
		Parameter:   params,
		Ctx:         ctx,
		StreamWrite: streamWrite,
	}

	if err := ctl.appService.GetRequestDataStream(cmd); err != nil {
		ctx.SSEvent("error", err.Error())
	}

}

// @Summary  CanRead
// @Description  check permission for read space app
// @Tags     SpaceAppWeb
// @Param    owner  path  string  true  "owner of space" MaxLength(40)
// @Param    name   path  string  true  "name of space" MaxLength(100)
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData
// @x-example {"data": "successfully"}
// @Router   /v1/space-app/{owner}/{name}/read [get]
func (ctl *SpaceAppWebController) CanRead(ctx *gin.Context) {
	index, err := ctl.parseIndex(ctx)
	if err != nil {
		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	if err := ctl.appService.CheckPermissionRead(user, &index); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, "successfully")
	}
}
