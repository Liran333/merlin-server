package repositoryadapter

import (
	"gorm.io/gorm"

	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
)

var appRepositoryAdapterInstance *appRepositoryAdapter

func Init(db *gorm.DB, tables *Tables) error {
	// must set branchTableName before migrating
	spaceappTableName = tables.SpaceApp

	if err := db.AutoMigrate(&spaceappDO{}); err != nil {
		return err
	}

	appRepositoryAdapterInstance = &appRepositoryAdapter{
		dao: postgresql.DAO(tables.SpaceApp),
	}

	return nil
}

func AppRepositoryAdapter() *appRepositoryAdapter {
	return appRepositoryAdapterInstance
}
