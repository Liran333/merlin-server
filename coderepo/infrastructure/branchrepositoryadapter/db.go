package branchrepositoryadapter

import (
	"gorm.io/gorm"

	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
)

var (
	branchAdapterInstance *branchAdapter
)

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

func BranchAdapter() *branchAdapter {
	return branchAdapterInstance
}
