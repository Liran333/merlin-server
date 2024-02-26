/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package domain provides domain models and types for the code repository branch.
package domain

import (
	coderepoprimitive "github.com/openmerlin/merlin-server/coderepo/domain/primitive"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

// Branch represents a code repository branch.
type Branch struct {
	BranchIndex

	Id         primitive.Identity
	RepoType   coderepoprimitive.RepoType
	CreatedAt  int64
	BaseBranch coderepoprimitive.BranchName
}

// BranchIndex represents the index information of a code repository branch.
type BranchIndex struct {
	Repo   primitive.MSDName
	Owner  primitive.Account
	Branch coderepoprimitive.BranchName
}

// RepoIndex returns the code repository index based on the branch index.
func (index *BranchIndex) RepoIndex() CodeRepoIndex {
	return CodeRepoIndex{
		Owner: index.Owner,
		Name:  index.Repo,
	}
}
