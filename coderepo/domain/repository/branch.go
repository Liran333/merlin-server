/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repository provides adapters for interacting with branches in a code repository.
package repository

import (
	"context"

	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

// BranchRepositoryAdapter represents an interface for managing branches in a code repository.
type BranchRepositoryAdapter interface {
	Add(*domain.Branch) error
	Delete(context.Context, primitive.Identity) error
	FindByIndex(context.Context, *domain.BranchIndex) (domain.Branch, error)
}

// BranchClientAdapter represents an interface for interacting with branches in a code repository client.
type BranchClientAdapter interface {
	CreateBranch(*domain.Branch) (string, error)
	DeleteBranch(*domain.BranchIndex) error
}
