/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package branchrepositoryadapter

import (
	"gorm.io/gorm"

	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
)

var (
	branchAdapterInstance *branchAdapter
)

// Init initializes the branch module by performing necessary setup and migrations.
func Init(db *gorm.DB, tables *Tables) error {
	// must set branchTableName before migrating
	branchTableName = tables.Branch

	if err := db.AutoMigrate(&branchDO{}); err != nil {
		return err
	}

	dao := postgresql.DAO(tables.Branch)

	branchAdapterInstance = &branchAdapter{dao: dao}

	return nil
}

// BranchAdapter returns an instance of the branchAdapter.
func BranchAdapter() *branchAdapter {
	return branchAdapterInstance
}
