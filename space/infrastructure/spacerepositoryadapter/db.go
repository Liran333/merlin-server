/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package spacerepositoryadapter

import "gorm.io/gorm"

var (
	spaceAdapterInstance       *spaceAdapter
	spaceLabelsAdapterInstance *spaceLabelsAdapter
	spaceModelInstance         *modelSpaceRelationAdapter
)

// Init initializes the database and sets up the necessary adapters.
func Init(db *gorm.DB, tables *Tables) error {
	// must set spaceTableName before migrating
	spaceTableName = tables.Space
	spaceModelRelationTableName = tables.SpaceModel

	if err := db.AutoMigrate(&spaceDO{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&modelSpaceRelationDO{}); err != nil {
		return err
	}

	dbInstance = db

	spaceDao := daoImpl{table: tables.Space}
	spaceModelDao := daoImpl{table: tables.SpaceModel}

	spaceAdapterInstance = &spaceAdapter{daoImpl: spaceDao}
	spaceLabelsAdapterInstance = &spaceLabelsAdapter{daoImpl: spaceDao}
	spaceModelInstance = &modelSpaceRelationAdapter{daoImpl: spaceModelDao}

	return nil
}

// SpaceAdapter returns the instance of the space adapter.
func SpaceAdapter() *spaceAdapter {
	return spaceAdapterInstance
}

// SpaceLabelsAdapter returns the instance of the space labels adapter.
func SpaceLabelsAdapter() *spaceLabelsAdapter {
	return spaceLabelsAdapterInstance
}

func ModelSpaceRelationAdapter() *modelSpaceRelationAdapter {
	return spaceModelInstance
}
