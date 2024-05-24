/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package activityrepositoryadapter provides an adapter for the model repository
package activityrepositoryadapter

import "gorm.io/gorm"

var (
	activityAdapterInstance *activityAdapter
)

// Init initializes the activity module by performing necessary setup and migrations.
func Init(db *gorm.DB, tables *Tables) error {
	// must set modelTableName before migrating
	activiyTableName = tables.Activity

	if err := db.AutoMigrate(&activityDO{}); err != nil {
		return err
	}

	dbInstance = db

	dao := daoImpl{table: tables.Activity}

	activityAdapterInstance = &activityAdapter{daoImpl: dao}

	return nil
}

// ActivityAdapter returns the instance of modelAdapter.
func ActivityAdapter() *activityAdapter {
	return activityAdapterInstance
}
