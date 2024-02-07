package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/coderepo/app"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

func ToCmdToFile(ctx *gin.Context) (*app.CmdToFile, error) {
	owner, err := primitive.NewAccount(ctx.Param("owner"))
	if err != nil {
		return nil, err
	}

	name, err := primitive.NewMSDName(ctx.Param("name"))
	if err != nil {
		return nil, err
	}

	fileRef := ctx.DefaultQuery("ref", "master")
	ref, err := primitive.NewCodeFileRef(fileRef)
	if err != nil {
		return nil, err
	}

	filePath := ctx.DefaultQuery("path", "/")
	path, err := primitive.NewCodeFilePath(filePath)
	if err != nil {
		return nil, err
	}

	return &app.CmdToFile{Owner: owner, Name: name, Ref: ref, FilePath: path}, nil
}
