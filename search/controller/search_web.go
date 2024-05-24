/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides functionality for managing the application's controllers.
package controller

import (
	"github.com/gin-gonic/gin"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/search/app"
)

// AddRouteForSearchWebController add route for search web
func AddRouteForSearchWebController(
	r *gin.RouterGroup,
	s app.SearchAppService,
	l middleware.OperationLog,
	m middleware.UserMiddleWare,
	rl middleware.RateLimiter,
) {

	ctl := &SearchWebController{}
	ctl.searchApp = s
	ctl.m = m

	r.GET("/v1/search", m.Optional, rl.CheckLimit, ctl.Search)
}

// SearchWebController search web controller
type SearchWebController struct {
	searchApp app.SearchAppService
	m         middleware.UserMiddleWare
}

// @Summary  List
// @Description  get model and space and org and user
// @Tags     SearchWeb
// @Param    searchKey     query  string  true "filter by name" MaxLength(100)
// @Param    type  query  []string  true "type of space/model/dataset/org/user" Enums(space,model,dataset,org,user)
// @Param  	 size  query  int    false "page data size" Maximum(100)
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=app.SearchDTO,msg=string,code=string}
// @Router /v1/search [get]
func (ctl *SearchWebController) Search(ctx *gin.Context) {
	var req quickSearchRequest
	if err := ctx.BindQuery(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
		return
	}

	user := ctl.m.GetUser(ctx)

	dto, err := ctl.searchApp.Search(&cmd, user)
	if err != nil {
		commonctl.SendError(ctx, err)
		return
	}

	commonctl.SendRespOfGet(ctx, &dto.ResultSet)
}
