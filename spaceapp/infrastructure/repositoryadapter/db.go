/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package repositoryadapter

import (
	"gorm.io/gorm"

	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
)

var (
	buildLogAdapterInstance      *buildLogAdapterImpl
	appRepositoryAdapterInstance *appRepositoryAdapter
)

// Init initializes the space app module by performing necessary setup and migrations.
func Init(db *gorm.DB, tables *Tables) error {
	// must set branchTableName before migrating
	spaceappTableName = tables.SpaceApp

	if err := db.AutoMigrate(&spaceappDO{}); err != nil {
		return err
	}

	dao := postgresql.DAO(tables.SpaceApp)

	appRepositoryAdapterInstance = &appRepositoryAdapter{
		dao: dao,
	}

	buildLogAdapterInstance = &buildLogAdapterImpl{
		dao: dao,
	}

	return nil
}

// AppRepositoryAdapter is an instance of the AppRepositoryAdapter.
func AppRepositoryAdapter() *appRepositoryAdapter {
	return appRepositoryAdapterInstance
}

func BuildLogAdapter() *buildLogAdapterImpl {
	return buildLogAdapterInstance
}
