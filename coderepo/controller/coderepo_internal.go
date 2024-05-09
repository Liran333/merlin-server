/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides the controllers for handling restful requests and converting them into commands
package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/openmerlin/merlin-server/coderepo/domain/resourceadapter"
	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	modelapp "github.com/openmerlin/merlin-server/models/app"
	spaceapp "github.com/openmerlin/merlin-server/space/app"
)

// AddRouteForCodeRepoStatisticInternalController adds routes for StatisticInternalController.
func AddRouteForCodeRepoStatisticInternalController(
	r *gin.RouterGroup,
	a resourceadapter.ResourceAdapter,
	m middleware.UserMiddleWare,
	o modelapp.ModelInternalAppService,
	p spaceapp.SpaceInternalAppService,
) {

	ctl := StatisticInternalController{
		repo:             a,
		modelInternalApp: o,
		spaceInternalApp: p,
	}

	r.PUT(`/v1/coderepo/:id/statistic`, m.Write, ctl.Update)
}

// StatisticInternalController is a struct that holds the necessary services
// and adapters for handling statistical operations.
type StatisticInternalController struct {
	repo             resourceadapter.ResourceAdapter
	modelInternalApp modelapp.ModelInternalAppService
	spaceInternalApp spaceapp.SpaceInternalAppService
}

// @Summary  Update
// @Description  update the download count of a model/space
// @Tags     Statistic
// @Param    id    path  string   true  "id of model/space" MaxLength(20)
// @Param    body  body  repoStatistics  true  "body of updating model/space info"
// @Accept   json
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Security Internal
// @Router   /v1/coderepo/{id}/statistic [put]
func (ctl *StatisticInternalController) Update(ctx *gin.Context) {
	repoId, err := primitive.NewIdentity(ctx.Param("id"))
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}

	var stats repoStatistics
	if err := ctx.BindJSON(&stats); err != nil {
		commonctl.SendBadRequestBody(ctx, err)
		return
	}

	repo, err := ctl.repo.GetByIndex(repoId)
	if err != nil {
		commonctl.SendError(ctx, err)
		return
	}

	switch repo.ResourceType() {
	case primitive.ObjTypeModel:
		middleware.SetAction(ctx, fmt.Sprintf("Update model statistics, ID: %v", repoId))
		err = ctl.modelInternalApp.UpdateStatistics(repoId,
			&modelapp.CmdToUpdateStatistics{DownloadCount: stats.DownloadCount})
	case primitive.ObjTypeSpace:
		middleware.SetAction(ctx, fmt.Sprintf("Update space statistics, ID: %v", repoId))
		err = ctl.spaceInternalApp.UpdateStatistics(repoId,
			&spaceapp.CmdToUpdateStatistics{DownloadCount: stats.DownloadCount})
	default:
		commonctl.SendError(ctx, fmt.Errorf("unsupported resource type"))
		return
	}

	if err != nil {
		commonctl.SendError(ctx, err)
		return
	}

	commonctl.SendRespOfPut(ctx, nil)
}
