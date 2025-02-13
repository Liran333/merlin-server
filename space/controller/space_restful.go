/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"github.com/gin-gonic/gin"

	activityapp "github.com/openmerlin/merlin-server/activity/app"
	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/space/app"
	userapp "github.com/openmerlin/merlin-server/user/app"
)

// AddRouteForSpaceRestfulController adds routes to the given router group for the SpaceRestfulController.
func AddRouteForSpaceRestfulController(
	r *gin.RouterGroup,
	s app.SpaceAppService,
	m middleware.UserMiddleWare,
	l middleware.OperationLog,
	sl middleware.SecurityLog,
	rl middleware.RateLimiter,
	u userapp.UserService,
	p middleware.PrivacyCheck,
	a activityapp.ActivityAppService,
) {
	ctl := SpaceRestfulController{
		SpaceController: SpaceController{
			appService:          s,
			userMiddleWare:      m,
			rateLimitMiddleWare: rl,
			user:                u,
			activity:            a,
		},
	}

	addRouteForSpaceController(r, &ctl.SpaceController, l, sl, rl)

	r.GET("/v1/space/:owner/:name", p.CheckOwner, m.Optional, rl.CheckLimit, ctl.Get)
	r.GET("/v1/space", m.Optional, rl.CheckLimit, ctl.List)
	r.POST("/v1/space/report", m.Optional, rl.CheckLimit, ctl.SpaceReport)
}

// SpaceRestfulController is a struct that holds the necessary dependencies for handling space-related operations.
type SpaceRestfulController struct {
	SpaceController
}

// @Summary  Get
// @Description  get space
// @Tags     SpaceRestful
// @Param    owner  path  string  true  "owner of space" MaxLength(40)
// @Param    name   path  string  true  "name of space" MaxLength(100)
// @Accept   json
// @Security Bearer
// @Success  200  {object}  commonctl.ResponseData{data=spaceDetail,msg=string,code=string}
// @Router   /v1/space/{owner}/{name} [get]
func (ctl *SpaceRestfulController) Get(ctx *gin.Context) {
	index, err := ctl.parseIndex(ctx)
	if err != nil {
		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	dto, err := ctl.appService.GetByName(ctx.Request.Context(), user, &index)
	if err != nil {
		commonctl.SendError(ctx, err)
		return
	}

	liked := false

	spaceId, _ := primitive.NewIdentity(dto.Id)
	if user != nil {
		liked, err = ctl.activity.HasLike(user, spaceId)
		if err != nil {
			commonctl.SendError(ctx, err)
			return
		}
	}

	detail := spaceDetail{
		Liked:    liked,
		SpaceDTO: &dto,
	}

	if avatar, err := ctl.user.GetUserAvatarId(ctx.Request.Context(), index.Owner); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		detail.OwnerAvatarId = avatar.AvatarId

		commonctl.SendRespOfGet(ctx, &detail)
	}
}

// @Summary  List
// @Description  list global public space
// @Tags     SpaceRestful
// @Param    name            query  string  false  "name of space" MaxLength(100)
// @Param    task            query  string  false  "task label" MaxLength(100)
// @Param    owner           query  string  true   "owner of space" MaxLength(40)
// @Param    license         query  string  false  "license label" MaxLength(40)
// @Param    framework       query  string  false  "framework " Enums(pytorch, mindspore)
// @Param    hardware_type   query  string  false  "type of space" Enums(npu, cpu)
// @Param    count           query  bool    false  "whether to calculate the total" Enums(true, false)
// @Param    has_app_file    query  bool    false  "show has app file space" Enums(true, false)
// @Param    sort_by         query  string  false  "sort types: most_likes, alphabetical, most_downloads, recently_updated, recently_created" Enums(most_likes, alphabetical,most_downloads,recently_updated,recently_created)
// @Param    page_num        query  int     false  "page num which starts from 1" Mininum(1)
// @Param    count_per_page  query  int     false  "count per page" MaxCountPerPage(100)
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=app.SpacesDTO,msg=string,code=string}
// @Router   /v1/space [get]
func (ctl *SpaceRestfulController) List(ctx *gin.Context) {
	var req restfulReqToListSpaces

	if err := ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	if result, err := ctl.appService.List(ctx.Request.Context(), user, &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, result)
	}
}

// @Summary  SendReportMail
// @Description  send report Email
// @Tags     SpaceRestful
// @Param    body  body      reqReportSpaceEmail  true  "body of send report"
// @Accept   json
// @Security Bearer
// @Success  201   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Router   /v1/space/report [post]
func (ctl *SpaceRestfulController) SpaceReport(ctx *gin.Context) {
	var req reqReportSpaceEmail
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, err)
		return
	}
	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}
	user := ctl.userMiddleWare.GetUser(ctx)
	if err := ctl.appService.SendSpaceReportEmail(user, &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, "succes send Email")
	}
}
