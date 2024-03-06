/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/space/app"
	userapp "github.com/openmerlin/merlin-server/user/app"
)

// AddRouteForSpaceWebController adds routes to the given router group for the SpaceWebController.
func AddRouteForSpaceWebController(
	r *gin.RouterGroup,
	s app.SpaceAppService,
	m middleware.UserMiddleWare,
	l middleware.OperationLog,
	rl middleware.RateLimiter,
	u userapp.UserService,
) {
	ctl := SpaceWebController{
		SpaceController: SpaceController{
			appService:     s,
			userMiddleWare: m,
			rateLimitMiddleWare: rl,
			user:           u,
		},
	}

	addRouteForSpaceController(r, &ctl.SpaceController, l, rl)

	r.GET("/v1/space/:owner/:name", m.Optional, rl.CheckLimit, ctl.Get)
	r.GET("/v1/space/:owner", m.Optional, rl.CheckLimit, ctl.List)
	r.GET("/v1/space", m.Optional, rl.CheckLimit, ctl.ListGlobal)
}

// SpaceWebController is a struct that holds the necessary dependencies for handling space-related operations in web controller.
type SpaceWebController struct {
	SpaceController
}

// @Summary  Get
// @Description  get space
// @Tags     SpaceWeb
// @Param    owner  path  string  true  "owner of space"
// @Param    name   path  string  true  "name of space"
// @Accept   json
// @Success  200  {object}  spaceDetail
// @Router   /v1/space/{owner}/{name} [get]
func (ctl *SpaceWebController) Get(ctx *gin.Context) {
	index, err := ctl.parseIndex(ctx)
	if err != nil {
		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	dto, err := ctl.appService.GetByName(user, &index)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	detail := spaceDetail{
		Liked:    true,
		SpaceDTO: &dto,
	}

	if user != nil {
		//TODO check whether user like the space
	}

	if avatar, err := ctl.user.GetUserAvatarId(index.Owner); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		detail.AvatarId = avatar.AvatarId

		commonctl.SendRespOfGet(ctx, &detail)
	}
}

// @Summary  List
// @Description  list space
// @Tags     SpaceWeb
// @Param    owner           path   string  true   "owner of space"
// @Param    name            query  string  false  "name of space"
// @Param    count           query  bool    false  "whether to calculate the total"
// @Param    sort_by         query  string  false  "sort types: most_likes, alphabetical, most_downloads, recently_updated, recently_created"
// @Param    page_num        query  int     false  "page num which starts from 1"
// @Param    count_per_page  query  int     false  "count per page"
// @Accept   json
// @Success  200  {object}  userSpacesInfo
// @Router   /v1/space/{owner} [get]
func (ctl *SpaceWebController) List(ctx *gin.Context) {
	var req reqToListUserSpaces
	if err := ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	cmd.Owner, err = primitive.NewAccount(ctx.Param("owner"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	dto, err := ctl.appService.List(user, &cmd)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	result := userSpacesInfo{
		Owner:     cmd.Owner.Account(),
		SpacesDTO: &dto,
	}

	if avatar, err := ctl.user.GetUserAvatarId(cmd.Owner); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		result.AvatarId = avatar.AvatarId

		commonctl.SendRespOfGet(ctx, &result)
	}
}

// @Summary  ListGlobal
// @Description  list global public space
// @Tags     SpaceWeb
// @Param    name            query  string  false  "name of space"
// @Param    task            query  string  false  "task label"
// @Param    others          query  string  false  "other labels, separate multiple each ones with commas"
// @Param    license         query  string  false  "license label"
// @Param    frameworks      query  string  false  "framework labels, separate multiple each ones with commas"
// @Param    count           query  bool    false  "whether to calculate the total"
// @Param    sort_by         query  string  false  "sort types: most_likes, alphabetical, most_downloads, recently_updated, recently_created"
// @Param    page_num        query  int     false  "page num which starts from 1"
// @Param    count_per_page  query  int     false  "count per page"
// @Accept   json
// @Success  200  {object}  spacesInfo
// @Router   /v1/space [get]
func (ctl *SpaceWebController) ListGlobal(ctx *gin.Context) {
	var req reqToListGlobalSpaces
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

	result, err := ctl.appService.List(user, &cmd)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	if v, err := ctl.setAvatars(&result); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, v)
	}
}

func (ctl *SpaceWebController) setAvatars(dto *app.SpacesDTO) (spacesInfo, error) {
	ms := dto.Spaces

	// get avatars

	v := map[string]bool{}
	for i := range ms {
		v[ms[i].Owner] = true
	}

	accounts := make([]primitive.Account, len(v))

	i := 0
	for k := range v {
		accounts[i] = primitive.CreateAccount(k)
		i++
	}

	avatars, err := ctl.user.GetUsersAvatarId(accounts)
	if err != nil {
		return spacesInfo{}, err
	}

	// set avatars

	am := map[string]string{}
	for i := range avatars {
		item := &avatars[i]

		am[item.Name] = item.AvatarId
	}

	infos := make([]spaceInfo, len(ms))
	for i := range ms {
		item := &ms[i]

		infos[i] = spaceInfo{
			AvatarId:     am[item.Owner],
			SpaceSummary: item,
		}
	}

	return spacesInfo{
		Total:  dto.Total,
		Spaces: infos,
	}, nil
}
