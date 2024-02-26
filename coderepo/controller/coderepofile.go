/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/coderepo/app"
	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/common/controller/middleware"
)

// AddRouterForCodeRepoFileController adds routes for CodeRepoFileController to the given router group.
func AddRouterForCodeRepoFileController(
	rg *gin.RouterGroup,
	cr app.CodeRepoFileAppService,
	m middleware.UserMiddleWare) {
	ctl := CodeRepoFileController{
		codeRepoFile:   cr,
		userMiddleWare: m,
	}

	rg.GET("/v1/files/:owner/:name", ctl.List)
	rg.GET("/v1/file/:owner/:name", ctl.Get)
	rg.GET("/v1/file/download/:owner/:name", ctl.Download)

}

// CodeRepoFileController is a struct that holds code repo file and user middleware for file operations.
type CodeRepoFileController struct {
	codeRepoFile   app.CodeRepoFileAppService
	userMiddleWare middleware.UserMiddleWare
}

// List handles the request to list files in a repository.
func (ctl *CodeRepoFileController) List(ctx *gin.Context) {
	codeRepoFile, err := ToCmdToFile(ctx)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
	}
	data, err := ctl.codeRepoFile.List(codeRepoFile)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, data)
	}
}

// Get handles the request to get a specific file in a repository.
func (ctl *CodeRepoFileController) Get(ctx *gin.Context) {
	codeRepoFile, err := ToCmdToFile(ctx)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
	}

	data, err := ctl.codeRepoFile.Get(codeRepoFile)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, data)
	}

}

// Download handles the request to download a specific file from a repository.
func (ctl *CodeRepoFileController) Download(ctx *gin.Context) {
	codeRepoFile, err := ToCmdToFile(ctx)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
	}

	data, err := ctl.codeRepoFile.Download(codeRepoFile)
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, data)
	}

}
