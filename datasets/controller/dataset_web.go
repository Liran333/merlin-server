/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides functionality for managing the application's controllers.
package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"

	activityapp "github.com/openmerlin/merlin-server/activity/app"
	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/datasets/app"
	userapp "github.com/openmerlin/merlin-server/user/app"
)

// AddRouteForDatasetWebController adds a router for the DatasetWebController with the given middleware.
func AddRouteForDatasetWebController(
	r *gin.RouterGroup,
	s app.DatasetAppService,
	m middleware.UserMiddleWare,
	l middleware.OperationLog,
	sl middleware.SecurityLog,
	u userapp.UserService,
	rl middleware.RateLimiter,
	p middleware.PrivacyCheck,
	a activityapp.ActivityAppService,
) {
	ctl := DatasetWebController{
		DatasetController: DatasetController{
			appService:     s,
			userMiddleWare: m,
			user:           u,
			activity:       a,
		},
	}

	addRouteForDatasetController(r, &ctl.DatasetController, l, sl)

	r.GET("/v1/dataset/:owner/:name", p.CheckOwner, m.Optional, ctl.Get)
	r.GET("/v1/dataset/:owner", p.CheckOwner, m.Optional, ctl.List)
	r.GET("/v1/dataset", m.Optional, ctl.ListGlobal)

	r.PUT("/v1/dataset/:id/disable", ctl.DatasetController.userMiddleWare.Write, l.Write, ctl.disable)
	r.POST("/v1/dataset/web/report", rl.CheckLimit, m.Write, l.Write, ctl.Report)
}

// DatasetWebController is a struct that holds the app service for dataset web operations.
type DatasetWebController struct {
	DatasetController
}

// @Summary  Get
// @Description  get dataset
// @Tags     DatasetWeb
// @Param    owner  path  string  true  "owner of dataset" MaxLength(40)
// @Param    name   path  string  true  "name of dataset" MaxLength(100)
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=datasetDetail,msg=string,code=string}
// @Router   /v1/dataset/{owner}/{name} [get]
func (ctl *DatasetWebController) Get(ctx *gin.Context) {
	index, err := ctl.parseIndex(ctx)
	if err != nil {
		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	dto, err := ctl.appService.GetByName(ctx.Request.Context(), user, &index)
	if err != nil {
		commonctl.SendError(ctx, xerrors.Errorf("failed to get dataset, %w", err))

		return
	}

	liked := false

	datasetId, _ := primitive.NewIdentity(dto.Id)
	if user != nil {
		liked, err = ctl.activity.HasLike(user, datasetId)
		if err != nil {
			commonctl.SendError(ctx, xerrors.Errorf("%w", err))
			return
		}
	}

	detail := datasetDetail{
		Liked:      liked,
		DatasetDTO: &dto,
	}

	if avatar, err := ctl.user.GetUserAvatarId(ctx.Request.Context(), index.Owner); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		detail.AvatarId = avatar.AvatarId

		commonctl.SendRespOfGet(ctx, &detail)
	}
}

// @Summary  List
// @Description  list dataset
// @Tags     DatasetWeb
// @Param    owner           path   string  true   "owner of dataset" MaxLength(40)
// @Param    name            query  string  false  "name of dataset" MaxLength(100)
// @Param    count           query  bool    false  "whether to calculate the total" Enums(true, false)
// @Param    sort_by         query  string  false  "sort types: most_likes, alphabetical, most_downloads, recently_updated, recently_created" Enums(most_likes, alphabetical,most_downloads,recently_updated,recently_created)
// @Param    page_num        query  int     false  "page num which starts from 1" Mininum(1)
// @Param    count_per_page  query  int     false  "count per page" MaxCountPerPage(100)
// @Accept   json
// @Success  200  {object}  userDatasetsInfo
// @Router   /v1/dataset/{owner} [get]
func (ctl *DatasetWebController) List(ctx *gin.Context) {
	var req reqToListUserDatasets
	if err := ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to parse req, %w", err))

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to convert req to cmd, %w", err))

		return
	}

	cmd.Owner, err = primitive.NewAccount(ctx.Param("owner"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("%w", err))

		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	dto, err := ctl.appService.List(ctx.Request.Context(), user, &cmd)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	result := userDatasetsInfo{
		Owner:       cmd.Owner.Account(),
		DatasetsDTO: &dto,
	}

	if userInfo, err := ctl.user.GetOrgOrUser(ctx.Request.Context(), user, cmd.Owner); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		result.AvatarId = userInfo.AvatarId
		result.OwnerType = userInfo.Type

		commonctl.SendRespOfGet(ctx, &result)
	}
}

// @Summary  ListGlobal
// @Description  list global public dataset
// @Tags     DatasetWeb
// @Param    name            query  string  false  "name of dataset" MaxLength(100)
// @Param    task            query  string  false  "task labels, separate multiple each ones with commas" MaxLength(100)
// @Param    license         query  string  false  "license label" MaxLength(40)
// @Param    size            query  string  false  "size labels" MaxLength(40)
// @Param    language        query  string  false  "language labels, separate multiple each ones with commas" MaxLength(100)
// @Param    domain          query  string  false  "domain labels, separate multiple each ones with commas" MaxLength(100)
// @Param    count           query  bool    false  "whether to calculate the total" Enums(true, false)
// @Param    sort_by         query  string  false  "sort types: most_likes, alphabetical, most_downloads, recently_updated, recently_created" Enums(most_likes, alphabetical,most_downloads,recently_updated,recently_created)
// @Param    page_num        query  int     false  "page num which starts from 1" Mininum(1)
// @Param    count_per_page  query  int     false  "count per page" MaxCountPerPage(100)
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=datasetsInfo,msg=string,code=string}
// @Router   /v1/dataset [get]
func (ctl *DatasetWebController) ListGlobal(ctx *gin.Context) {
	var req reqToListGlobalDatasets
	if err := ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to parse req, %w", err))

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to convert req to cmd, %w", err))

		return
	}

	user := ctl.userMiddleWare.GetUser(ctx)

	result, err := ctl.appService.List(ctx.Request.Context(), user, &cmd)
	if err != nil {
		commonctl.SendError(ctx, err)

		return
	}

	if v, err := ctl.setUserInfo(ctx, &result); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, v)
	}
}

func (ctl *DatasetWebController) setUserInfo(ctx *gin.Context, dto *app.DatasetsDTO) (datasetsInfo, error) {
	ms := dto.Datasets

	// get avatars
	v := map[string]userapp.UserDTO{}
	for i := range ms {
		v[ms[i].Owner] = userapp.UserDTO{}
	}

	accounts := make([]primitive.Account, len(v))
	i := 0
	for k := range v {
		accounts[i] = primitive.CreateAccount(k)
		userInfo, err := ctl.user.GetOrgOrUser(ctx.Request.Context(), nil, accounts[i])
		if err != nil {
			return datasetsInfo{}, xerrors.Errorf("failed to get user info, %w", err)
		}
		v[k] = userInfo
		i++
	}

	// set avatars
	infos := make([]datasetInfo, len(ms))
	for i := range ms {
		item := &ms[i]

		infos[i] = datasetInfo{
			AvatarId:       v[item.Owner].AvatarId,
			OwnerType:      v[item.Owner].Type,
			Owner:          item.Owner,
			DatasetSummary: item,
		}
	}

	return datasetsInfo{
		Total:    dto.Total,
		Datasets: infos,
	}, nil
}

// @Summary  disable dataset
// @Description  disable the dataset
// @Tags     DatasetWeb
// @Param    id      path  string  true  "id of dataset" MaxLength(20)
// @Param    body  body      reqToDisableDataset  true  "body of disable dataset"
// @Accept   json
// @Security Bearer
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Router   /v1/dataset/{id}/disable [put]
func (ctl *DatasetWebController) disable(ctx *gin.Context) {
	middleware.SetAction(ctx, fmt.Sprintf("disable dataset of %s", ctx.Param("id")))

	req := reqToDisableDataset{}
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, xerrors.Errorf("failed to parse req, %w", err))

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to convert req to cmd, %w", err))

		return
	}

	datasetId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	action, err := ctl.appService.Disable(
		ctx.Request.Context(),
		ctl.userMiddleWare.GetUser(ctx),
		datasetId,
		&cmd,
	)

	middleware.SetAction(ctx, fmt.Sprintf("%s, set %s", action, req.action()))

	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPut(ctx, nil)
	}
}

// @Summary  SendReportMail
// @Description  send report Email
// @Tags     DatasetWeb
// @Param    body  body      reqReportDatasetEmail  true  "body of send report"
// @Accept   json
// @Security Bearer
// @Success  201   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Router   /v1/dataset/web/report [post]
func (ctl *DatasetWebController) Report(ctx *gin.Context) {
	var req reqReportDatasetEmail
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestBody(ctx, xerrors.Errorf("failed to parse req, %w", err))
		return
	}
	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestBody(ctx, err)
	}
	user := ctl.userMiddleWare.GetUser(ctx)
	if err := ctl.appService.SendReportMail(user, &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfPost(ctx, "success send Email")
	}
}
