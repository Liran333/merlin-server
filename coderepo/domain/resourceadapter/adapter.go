/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package resourceadapter provides adapters for retrieving resources from a code repository.
package resourceadapter

import (
	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/coderepo/domain/primitive"
)

// ResourceAdapter represents an interface for retrieving resources from a code repository.
type ResourceAdapter interface {
	GetByName(*domain.CodeRepoIndex) (domain.Resource, error)
	GetByIndex(identity primitive.Identity) (domain.Resource, error)
	GetByType(primitive.RepoType, *domain.CodeRepoIndex) (domain.Resource, error)
}
