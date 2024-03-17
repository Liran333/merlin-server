/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package controller

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/search/app"
)

type quickSearchRequest struct {
	SearchKey  string   `json:"search_key" binding:"required"`
	SearchType []string `json:"search_type" binding:"required"`
	Size       int      `json:"size" binding:"required"`
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
