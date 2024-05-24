/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package app provides application services for creating and managing branches.
package app

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/search/domain"
)

// CmdToSearch is the command to search
type CmdToSearch struct {
	SearchKey  primitive.SearchKey
	SearchType primitive.SearchType
	Size       primitive.Size
}

// SearchDTO is the data transfer object for search
type SearchDTO struct {
	ResultSet domain.SearchResult `json:"result_set"`
}
