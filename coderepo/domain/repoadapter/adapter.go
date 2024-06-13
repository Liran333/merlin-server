/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repoadapter provides interfaces for adapting code repository operations.
package repoadapter

import (
	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

// RepoAdapter is an interface that defines methods for code repository operations.
type RepoAdapter interface {
	Add(*domain.CodeRepo, bool) error
	Delete(*domain.CodeRepoIndex) error
	Save(*domain.CodeRepoIndex, *domain.CodeRepo) error
	FindByIndex(primitive.Identity) (domain.CodeRepo, error)
	IsNotFound(primitive.Identity) bool
}
