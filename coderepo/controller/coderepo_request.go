/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides the controllers for handling restful requests and converting them into commands
package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/coderepo/app"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

// ToCmdToCheckRepoExists is a function that parses the request context
// and returns a command to check if a repository exists.
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

// repoStatistics download count
type repoStatistics struct {
	DownloadCount int `json:"download_count"`
}

// repoVisitCount visit count
type repoVisitCount struct {
	VisitCount    int `json:"visit_count"`
}
