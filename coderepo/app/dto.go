/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"github.com/openmerlin/merlin-server/coderepo/domain"
	repoprimitive "github.com/openmerlin/merlin-server/coderepo/domain/primitive"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/utils"
)

// CmdToCreateRepo is a struct representing the command to create a repository.
type CmdToCreateRepo struct {
	Name       primitive.MSDName
	Owner      primitive.Account
	License    primitive.License
	Visibility primitive.Visibility
	InitReadme bool
}

func (cmd *CmdToCreateRepo) toCodeRepo(user primitive.Account) domain.CodeRepo {
	return domain.CodeRepo{
		Name:       cmd.Name,
		Owner:      cmd.Owner,
		License:    cmd.License,
		CreatedBy:  user,
		Visibility: cmd.Visibility,
	}
}

// CmdToUpdateRepo is a struct representing the command to update a repository.
type CmdToUpdateRepo struct {
	Name       primitive.MSDName
	Visibility primitive.Visibility
}

func (cmd *CmdToUpdateRepo) toRepo(repo *domain.CodeRepo) (b bool) {
	if v := cmd.Visibility; v != nil && v != repo.Visibility {
		repo.Visibility = v
		b = true
	}

	return
}

// CmdToCreateBranch is a struct representing the command to create a branch.
type CmdToCreateBranch struct {
	domain.BranchIndex

	RepoType   repoprimitive.RepoType
	BaseBranch repoprimitive.BranchName
}

func (cmd *CmdToCreateBranch) toBranch() domain.Branch {
	branch := domain.Branch{
		BranchIndex: cmd.BranchIndex,
		BaseBranch:  cmd.BaseBranch,
		RepoType:    cmd.RepoType,
		CreatedAt:   utils.Now(),
	}

	return branch
}

// BranchCreateDTO is a struct representing the data transfer object for creating a branch.
type BranchCreateDTO struct {
	Name string `json:"branch_name"`
}

func toBranchCreateDTO(v string) BranchCreateDTO {
	return BranchCreateDTO{
		Name: v,
	}
}

// CmdToDeleteBranch is a struct representing the command to delete a branch.
type CmdToDeleteBranch struct {
	domain.BranchIndex

	RepoType repoprimitive.RepoType
}

// CmdToCheckRepoExists is an alias type for the domain.CodeRepoIndex.
type CmdToCheckRepoExists = domain.CodeRepoIndex
