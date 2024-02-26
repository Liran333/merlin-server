/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package branchrepositoryadapter

import (
	"github.com/openmerlin/merlin-server/coderepo/domain"
	coderepoprimitive "github.com/openmerlin/merlin-server/coderepo/domain/primitive"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

var (
	branchTableName string
)

type branchDO struct {
	Id         int64  `gorm:"primaryKey;autoIncrement"`
	Repo       string `gorm:"column:repo;index:branch_index,unique,priority:2"`
	Owner      string `gorm:"column:owner;index:branch_index,unique,priority:1"`
	Branch     string `gorm:"column:branch;index:branch_index,unique,priority:3"`
	RepoType   string `gorm:"column:repo_type"`
	CreatedAt  int64  `gorm:"column:created_at"`
	BaseBranch string `gorm:"column:base_branch"`
}

func toBranchDO(m *domain.Branch) branchDO {
	return branchDO{
		Repo:       m.Repo.MSDName(),
		Owner:      m.Owner.Account(),
		Branch:     m.Branch.BranchName(),
		RepoType:   m.RepoType.RepoType(),
		CreatedAt:  m.CreatedAt,
		BaseBranch: m.BaseBranch.BranchName(),
	}
}

// TableName returns the table name for the branchDO struct.
func (do *branchDO) TableName() string {
	return branchTableName
}

func (do *branchDO) toBranch() domain.Branch {
	return domain.Branch{
		BranchIndex: domain.BranchIndex{
			Repo:   primitive.CreateMSDName(do.Repo),
			Owner:  primitive.CreateAccount(do.Owner),
			Branch: coderepoprimitive.CreateBranchName(do.RepoType),
		},
		Id:         primitive.CreateIdentity(do.Id),
		RepoType:   coderepoprimitive.CreateRepoType(do.RepoType),
		CreatedAt:  do.CreatedAt,
		BaseBranch: coderepoprimitive.CreateBranchName(do.BaseBranch),
	}
}
