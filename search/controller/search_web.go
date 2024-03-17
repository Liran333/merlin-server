/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/search/app"
)

const TypeModelResultNum = 6

func AddRouteForSearchWebController(
	r *gin.RouterGroup,
	s app.SearchAppService,
	l middleware.OperationLog,
	m middleware.UserMiddleWare,
) {

	ctl := &SearchWebController{}
	ctl.searchApp = s
	ctl.m = m

	r.GET("/v1/search", m.Optional, ctl.Search)
}

type SearchWebController struct {
	searchApp app.SearchAppService
	m         middleware.UserMiddleWare
}

func (ctl *SearchWebController) Search(ctx *gin.Context) {
	var req quickSearchRequest

	req.SearchKey = ctx.Query("searchKey")
	req.SearchType = ctx.QueryArray("type")
	req.Size, _ = strconv.Atoi(ctx.Query("size"))
	if req.Size == 0 {
		logrus.Infof("Failed to get size, set it to %d", TypeModelResultNum)
		req.Size = TypeModelResultNum
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
