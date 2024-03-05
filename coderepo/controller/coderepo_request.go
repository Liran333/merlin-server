package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/coderepo/app"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

func ToCmdToCheckRepoExists(ctx *gin.Context) (*app.CmdToCheckRepoExists, error) {
	owner, err := primitive.NewAccount(ctx.Param("owner"))
	if err != nil {
		return nil, err
	}

	name, err := primitive.NewMSDName(ctx.Param("name"))
	if err != nil {
		return nil, err
	}

	return &app.CmdToCheckRepoExists{
		Owner: owner,
		Name:  name,
	}, nil
}
