/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"

	activityapp "github.com/openmerlin/merlin-server/activity/app"
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
	ms app.ModelSpaceAppService,
	sv app.SpaceVariableService,
	ss app.SpaceSecretService,
	m middleware.UserMiddleWare,
	l middleware.OperationLog,
	sl middleware.SecurityLog,
	rl middleware.RateLimiter,
	u userapp.UserService,
	p middleware.PrivacyCheck,
	a activityapp.ActivityAppService,
) {
	ctl := SpaceWebController{
		SpaceController: SpaceController{
			appService:          s,
			variableService:     sv,
			secretService:       ss,
			userMiddleWare:      m,
			rateLimitMiddleWare: rl,
			user:                u,
			activity:            a,
		},
		modelSpaceService: ms,
	}

	addRouteForSpaceController(r, &ctl.SpaceController, l, sl, rl)

	addRouteForSpaceVariableController(r, &ctl.SpaceController, l, sl, rl)

	addRouteForSpaceSecretController(r, &ctl.SpaceController, l, sl, rl)

	r.GET("/v1/space/:owner/:name", p.CheckOwner, m.Optional, rl.CheckLimit, ctl.Get)
	r.GET("/v1/space/:owner", p.CheckOwner, m.Optional, rl.CheckLimit, ctl.List)
	r.GET("/v1/space", m.Optional, rl.CheckLimit, ctl.ListGlobal)
	r.GET("/v1/space/recommend", m.Optional, rl.CheckLimit, ctl.ListRecommends)
	r.GET("/v1/space/boutique", m.Optional, rl.CheckLimit, ctl.ListBoutiques)
	r.GET("/v1/space/relation/:id/model", m.Optional, rl.CheckLimit, ctl.GetModelsBySpaceId)

	r.PUT("/v1/space/:id/disable", ctl.SpaceController.userMiddleWare.Write, l.Write, rl.CheckLimit, ctl.Disable)
}

// SpaceWebController is a struct that holds the necessary dependencies for handling space-related operations in web controller.
type SpaceWebController struct {
	SpaceController
	modelSpaceService app.ModelSpaceAppService
}

// @Summary  Get
// @Description  get space
// @Tags     SpaceWeb
// @Param    owner  path  string  true  "owner of space" MaxLength(40)
// @Param    name   path  string  true  "name of space" MaxLength(100)
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=spaceDetail,msg=string,code=string}
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

	if userInfo, err := ctl.user.GetOrgOrUser(user, index.Owner); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		detail.OwnerAvatarId = userInfo.AvatarId
		detail.OwnerType = userInfo.Type

		commonctl.SendRespOfGet(ctx, &detail)
	}
}

// @Summary  List
// @Description  list space
// @Tags     SpaceWeb
// @Param    owner           path   string  true   "owner of space" MaxLength(40)
// @Param    name            query  string  false  "name of space" MaxLength(100)
// @Param    count           query  bool    false  "whether to calculate the total" Enums(true, false)
// @Param    sort_by         query  string  false  "sort types: most_likes, alphabetical, most_downloads, recently_updated, recently_created" Enums(most_likes, alphabetical,most_downloads,recently_updated,recently_created)
// @Param    page_num        query  int     false  "page num which starts from 1" Mininum(1)
// @Param    count_per_page  query  int     false  "count per page" MaxCountPerPage(100)
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=userSpacesInfo,msg=string,code=string}
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

	if userInfo, err := ctl.user.GetOrgOrUser(user, cmd.Owner); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		result.AvatarId = userInfo.AvatarId
		result.OwnerType = userInfo.Type

		commonctl.SendRespOfGet(ctx, &result)
	}
}

// @Summary  ListRecommends
// @Description  list recommend space
// @Tags     SpaceWeb
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=spacesRecommendInfo,msg=string,code=string}
// @Router   /v1/space/recommend [get]
func (ctl *SpaceWebController) ListRecommends(ctx *gin.Context) {
	user := ctl.userMiddleWare.GetUser(ctx)

	spacesDTO, err := ctl.appService.Recommend(user)
	if err != nil {
		commonctl.SendError(ctx, xerrors.Errorf("failed to get recommend spaces: %w", err))

		return
	}

	sps := make([]spaceRecommendInfo, 0, len(spacesDTO))

	for _, v := range spacesDTO {
		userInfo, err := ctl.user.GetOrgOrUser(nil, primitive.CreateAccount(v.Owner))
		if err != nil {
			commonctl.SendError(ctx, err)
		}

		s := v
		sp := spaceRecommendInfo{
			OwnerType: userInfo.Type,
			SpaceDTO:  &s,
		}

		sps = append(sps, sp)
	}

	result := spacesRecommendInfo{
		Spaces: sps,
	}

	commonctl.SendRespOfGet(ctx, &result)
}

// @Summary  ListBoutiques
// @Description  list boutique space
// @Tags     SpaceWeb
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=spacesRecommendInfo,msg=string,code=string}
// @Router   /v1/space/boutique [get]
func (ctl *SpaceWebController) ListBoutiques(ctx *gin.Context) {
	user := ctl.userMiddleWare.GetUser(ctx)

	spacesDTO, err := ctl.appService.Boutique(user)
	if err != nil {
		commonctl.SendError(ctx, xerrors.Errorf("failed to get boutique spaces: %w", err))

		return
	}

	sps := make([]spaceRecommendInfo, 0, len(spacesDTO))

	for _, v := range spacesDTO {
		s := v
		sp := spaceRecommendInfo{
			SpaceDTO: &s,
		}

		sps = append(sps, sp)
	}

	result := spacesRecommendInfo{
		Spaces: sps,
	}

	commonctl.SendRespOfGet(ctx, &result)
}

// @Summary  ListGlobal
// @Description  list global public space
// @Tags     SpaceWeb
// @Param    name            query  string  false  "name of space" MaxLength(100)
// @Param    domain          query  string  false  "domain label" MaxLength(100)
// @Param    license         query  string  false  "license label" MaxLength(40)
// @Param    framework       query  string  false  "framework " Enums(pytorch, mindspore)
// @Param    count           query  bool    false  "whether to calculate the total" Enums(true, false)
// @Param    sort_by         query  string  false  "sort types: most_likes, alphabetical, most_downloads, recently_updated, recently_created" Enums(most_likes, alphabetical,most_downloads,recently_updated,recently_created)
// @Param    page_num        query  int     false  "page num which starts from 1" Mininum(1)
// @Param    count_per_page  query  int     false  "count per page" MaxCountPerPage(100)
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
	v := map[string]userapp.UserDTO{}
	for i := range ms {
		v[ms[i].Owner] = userapp.UserDTO{}
	}

	accounts := make([]primitive.Account, len(v))

	i := 0
	for k := range v {
		accounts[i] = primitive.CreateAccount(k)
		userInfo, err := ctl.user.GetOrgOrUser(nil, accounts[i])
		if err != nil {
			return spacesInfo{}, err
		}
		v[k] = userInfo
		i++
	}

	// set avatars
	infos := make([]spaceInfo, len(ms))
	for i := range ms {
		item := &ms[i]

		infos[i] = spaceInfo{
			AvatarId:  v[item.Owner].AvatarId,
			OwnerType: v[item.Owner].Type,
			Owner:     item.Owner,

			SpaceSummary: item,
		}
	}

	return spacesInfo{
		Total:  dto.Total,
		Spaces: infos,
	}, nil
}

// @Summary  GetModelsBySpaceId
// @Description  get models related to a space
// @Tags     SpaceWeb
// @Param    id    path  string   true  "id of space" MaxLength(20)
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=[]app.SpaceModelDTO,msg=string,code=string}
// @Router   /v1/space/relation/{id}/model [get]
func (ctl *SpaceWebController) GetModelsBySpaceId(ctx *gin.Context) {
	spaceId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	models, err := ctl.modelSpaceService.GetModelsBySpaceId(user, spaceId)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	for _, model := range models {
		if avatar, err := ctl.user.GetUserAvatarId(primitive.CreateAccount(model.Owner)); err != nil {
			model.AvatarId = ""
		} else {
			model.AvatarId = avatar.AvatarId
		}
	}

	commonctl.SendRespOfGet(ctx, &models)
}
