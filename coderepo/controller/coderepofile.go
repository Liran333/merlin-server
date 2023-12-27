package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/openmerlin/merlin-server/common/controller/middleware"

	"github.com/openmerlin/merlin-server/coderepo/app"
	commonctl "github.com/openmerlin/merlin-server/common/controller"
)

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

type CodeRepoFileController struct {
	codeRepoFile   app.CodeRepoFileAppService
	userMiddleWare middleware.UserMiddleWare
}

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
