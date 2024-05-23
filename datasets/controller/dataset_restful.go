/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides functionality for managing the application's controllers.
package controller

import (
	"github.com/gin-gonic/gin"

	activityapp "github.com/openmerlin/merlin-server/activity/app"
	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/datasets/app"
	userapp "github.com/openmerlin/merlin-server/user/app"
)

// AddRouteForDatasetRestfulController adds a router for the DatasetRestfulController with the given middleware.
func AddRouteForDatasetRestfulController(
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
	ctl := DatasetRestfulController{
		DatasetController: DatasetController{
			appService:     s,
			userMiddleWare: m,
			user:           u,
			activity:       a,
		},
	}

	addRouteForDatasetController(r, &ctl.DatasetController, l, sl)

	r.GET("/v1/dataset/:owner/:name", p.CheckOwner, m.Optional, rl.CheckLimit, ctl.Get)
	r.GET("/v1/dataset", m.Optional, rl.CheckLimit, ctl.List)
}

// DatasetRestfulController is a struct that holds the app service for dataset restful operations.
type DatasetRestfulController struct {
	DatasetController
}

// @Summary  Get
// @Description  get dataset
// @Tags     DatasetRestful
// @Param    owner  path  string  true  "owner of dataset" MaxLength(40)
// @Param    name   path  string  true  "name of dataset" MaxLength(100)
// @Accept   json
// @Security Bearer
// @Success  200  {object}  commonctl.ResponseData{data=app.DatasetDTO,msg=string,code=string}
// @Router   /v1/dataset/{owner}/{name} [get]
func (ctl *DatasetRestfulController) Get(ctx *gin.Context) {
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

	datasetId, _ := primitive.NewIdentity(dto.Id)
	if user != nil {
		liked, err = ctl.activity.HasLike(user, datasetId)
		if err != nil {
			commonctl.SendError(ctx, err)
			return
		}
	}

	detail := datasetDetail{
		Liked:      liked,
		DatasetDTO: &dto,
	}

	if err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, &detail)
	}
}

// @Summary  List
// @Description  list global public datasets
// @Tags     DatasetRestful
// @Param    name            query  string  false  "name of dataset" MaxLength(100)
// @Param    task            query  string  false  "task labels, separate multiple each ones with commas" MaxLength(100)
// @Param    size            query  string  false  "size label" MaxLength(40)
// @Param    language        query  string  false  "language labels, separate multiple each ones with commas" MaxLength(100)
// @Param    domain          query  string  false  "domain labels, separate multiple each ones with commas" MaxLength(100)
// @Param    owner           query  string  true   "owner of dataset" MaxLength(40)
// @Param    license         query  string  false  "license label" MaxLength(40)
// @Param    count           query  bool    false  "whether to calculate the total" Enums(true, false)
// @Param    sort_by         query  string  false  "sort types: most_likes, alphabetical, most_downloads, recently_updated, recently_created" Enums(most_likes, alphabetical,most_downloads,recently_updated,recently_created)
// @Param    page_num        query  int     false  "page num which starts from 1" Mininum(1)
// @Param    count_per_page  query  int     false  "count per page" MaxCountPerPage(100)
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=app.DatasetsDTO,msg=string,code=string}
// @Router   /v1/dataset [get]
func (ctl *DatasetRestfulController) List(ctx *gin.Context) {
	var req restfulReqToListDatasets

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

	if result, err := ctl.appService.List(user, &cmd); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, result)
	}
}
