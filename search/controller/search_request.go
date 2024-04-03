/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides functionality for managing the application's controllers.
package controller

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/search/app"
)

type quickSearchRequest struct {
	SearchKey  string   `form:"searchKey" binding:"required"`
	SearchType []string `form:"type" binding:"required"`
	Size       int      `form:"size"`
}

func (req *quickSearchRequest) toCmd() (cmd app.CmdToSearch, err error) {
	if cmd.SearchKey, err = primitive.NewSearchKey(req.SearchKey); err != nil {
		return
	}

	if cmd.SearchType, err = primitive.NewSearchType(req.SearchType); err != nil {
		return
	}

	if cmd.Size, err = primitive.NewSize(req.Size); err != nil {
		return
	}

	return
}
