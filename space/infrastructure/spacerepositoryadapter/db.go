/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package spacerepositoryadapter

import "gorm.io/gorm"

var (
	spaceAdapterInstance       *spaceAdapter
	spaceLabelsAdapterInstance *spaceLabelsAdapter
)

// Init initializes the database and sets up the necessary adapters.
func Init(db *gorm.DB, tables *Tables) error {
	// must set spaceTableName before migrating
	spaceTableName = tables.Space

	if err := db.AutoMigrate(&spaceDO{}); err != nil {
		return err
	}

	dbInstance = db

	dao := daoImpl{table: tables.Space}

	spaceAdapterInstance = &spaceAdapter{daoImpl: dao}
	spaceLabelsAdapterInstance = &spaceLabelsAdapter{daoImpl: dao}

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
