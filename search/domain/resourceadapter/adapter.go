/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package resourceadapter provides adapter models and configuration for a specific functionality.
package resourceadapter

import (
	"context"

	"github.com/openmerlin/merlin-server/search/domain"
)

// ResourceAdapter is an interface for resource adapter
type ResourceAdapter interface {
	Search(ctx context.Context, opt *domain.SearchOption) (domain.SearchResult, error)
}
