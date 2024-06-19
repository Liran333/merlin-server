/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package branchrepositoryadapter provides an adapter for the branch repository using GORM.
package branchrepositoryadapter

import (
	"context"

	"gorm.io/gorm"

	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

type dao interface {
	DB() *gorm.DB
	GetRecord(ctx context.Context, filter, result interface{}) error
	DeleteByPrimaryKey(ctx context.Context, row interface{}) error
	EqualQuery(field string) string
}
type branchAdapter struct {
	dao
}

// Add adds a branch to the repository.
func (adapter *branchAdapter) Add(branch *domain.Branch) error {
	do := toBranchDO(branch)
	v := adapter.DB().Create(&do)

	return v.Error
}

// Delete deletes a branch from the repository by its ID.
func (adapter *branchAdapter) Delete(ctx context.Context, id primitive.Identity) error {
	return adapter.DeleteByPrimaryKey(
		ctx, &branchDO{Id: id.Integer()},
	)
}

// FindByIndex finds a branch in the repository by its index.
func (adapter *branchAdapter) FindByIndex(ctx context.Context, index *domain.BranchIndex) (domain.Branch, error) {
	do := branchDO{
		Owner:  index.Owner.Account(),
		Repo:   index.Repo.MSDName(),
		Branch: index.Branch.BranchName(),
	}
	if err := adapter.GetRecord(ctx, &do, &do); err != nil {
		return domain.Branch{}, err
	}

	return do.toBranch(), nil
}
