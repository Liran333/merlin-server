/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides the controllers for handling restful requests and converting them into commands
package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/coderepo/app"
	repoprimitive "github.com/openmerlin/merlin-server/coderepo/domain/primitive"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

type restfulReqToCreateBranch struct {
	Branch     string `json:"branch"       required:"true"`
	BaseBranch string `json:"base_branch"  required:"true"`
}

func (req *restfulReqToCreateBranch) action(ctx *gin.Context) string {
	return fmt.Sprintf("create branch %s from %s/%s/%s",
		req.Branch, ctx.Param("owner"), ctx.Param("repo"), req.BaseBranch)
}

func (req *restfulReqToCreateBranch) toCmd(ctx *gin.Context) (cmd app.CmdToCreateBranch, err error) {
	if cmd.Owner, err = primitive.NewAccount(ctx.Param("owner")); err != nil {
		return
	}
	if cmd.Repo, err = primitive.NewMSDName(ctx.Param("repo")); err != nil {
		return
	}
	if cmd.Branch, err = repoprimitive.NewBranchName(req.Branch); err != nil {
		return
	}
	if cmd.BaseBranch, err = repoprimitive.NewBranchName(req.BaseBranch); err != nil {
		return
	}
	if cmd.RepoType, err = repoprimitive.NewRepoType(ctx.Param("type")); err != nil {
		return
	}

	return
}

func toBanchDeleteCmd(ctx *gin.Context) (cmd app.CmdToDeleteBranch, err error) {
	if cmd.Owner, err = primitive.NewAccount(ctx.Param("owner")); err != nil {
		return
	}
	if cmd.Repo, err = primitive.NewMSDName(ctx.Param("repo")); err != nil {
		return
	}
	if cmd.RepoType, err = repoprimitive.NewRepoType(ctx.Param("type")); err != nil {
		return
	}
	if cmd.Branch, err = repoprimitive.NewBranchName(ctx.Param("branch")); err != nil {
		return
	}

	return
}
