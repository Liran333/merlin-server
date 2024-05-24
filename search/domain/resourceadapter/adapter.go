/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package resourceadapter provides adapter models and configuration for a specific functionality.
package resourceadapter

import "github.com/openmerlin/merlin-server/search/domain"

// ResourceAdapter is an interface for resource adapter
type ResourceAdapter interface {
	Search(opt *domain.SearchOption) (domain.SearchResult, error)
}
