package controller

import (
	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/space-app/app"
)

func AddRouteForSpaceAppInternalController(
	r *gin.RouterGroup,
	s app.SpaceAppInternalAppService,
	m middleware.UserMiddleWare,
) {

	ctl := SpaceAppInternalController{
		appService: s,
	}

	r.POST(`/v1/space-app`, m.Write, ctl.Create)
	r.PUT(`/v1/space-app/build/started`, m.Write, ctl.NotifyBuildIsStarted)
	r.PUT(`/v1/space-app/build/done`, m.Write, ctl.NotifyBuildIsDone)
	r.PUT(`/v1/space-app/service/started`, m.Write, ctl.NotifyServiceIsStarted)
}

type SpaceAppInternalController struct {
	appService app.SpaceAppInternalAppService
}

// @Summary  Create
// @Description  create space app
// @Tags     SpaceApp
// @Param    body  body  reqToCreateSpaceApp  true  "body of creating space app"
// @Accept   json
// @Success  201   {object}  commonctl.ResponseData
// @Security Bearer
// @Router   /v1/space-app/ [post]
func (ctl *SpaceAppInternalController) Create(ctx *gin.Context) {
	req := reqToCreateSpaceApp{}
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if err := ctl.appService.Create(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, "successfully")
	}
}

// @Summary  NotifyBuildIsStarted
// @Description  notidy space app building is started
// @Tags     SpaceApp
// @Param    body  body  reqToUpdateBuildInfo  true  "body"
// @Accept   json
// @Success  202   {object}  commonctl.ResponseData
// @Security Bearer
// @Router   /v1/space-app/build/started [put]
func (ctl *SpaceAppInternalController) NotifyBuildIsStarted(ctx *gin.Context) {
	req := reqToUpdateBuildInfo{}

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if err := ctl.appService.NotifyBuildIsStarted(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  NotifyBuildIsDone
// @Description  notidy space app build is done
// @Tags     SpaceApp
// @Param    body  body  reqToSetBuildIsDone  true  "body"
// @Accept   json
// @Success  202   {object}  commonctl.ResponseData
// @Security Bearer
// @Router   /v1/space-app/build/done [put]
func (ctl *SpaceAppInternalController) NotifyBuildIsDone(ctx *gin.Context) {
	req := reqToSetBuildIsDone{}

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if err := ctl.appService.NotifyBuildIsDone(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  NotifyServiceIsStarted
// @Description  notidy space app service is started
// @Tags     SpaceApp
// @Param    body  body  reqToUpdateServiceInfo  true  "body"
// @Accept   json
// @Success  202   {object}  commonctl.ResponseData
// @Security Bearer
// @Router   /v1/space-app/service/started [put]
func (ctl *SpaceAppInternalController) NotifyServiceIsStarted(ctx *gin.Context) {
	req := reqToUpdateServiceInfo{}

	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if err := ctl.appService.NotifyServiceIsStarted(&cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}
