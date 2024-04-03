/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package modelrepositoryadapter provides an adapter for the model repository
package modelrepositoryadapter

import "gorm.io/gorm"

var (
	modelAdapterInstance       *modelAdapter
	modelLabelsAdapterInstance *modelLabelsAdapter
)

// Init initializes the model module by performing necessary setup and migrations.
func Init(db *gorm.DB, tables *Tables) error {
	// must set modelTableName before migrating
	modelTableName = tables.Model

	if err := db.AutoMigrate(&modelDO{}); err != nil {
		return err
	}

	dbInstance = db

	dao := daoImpl{table: tables.Model}

	modelAdapterInstance = &modelAdapter{daoImpl: dao}
	modelLabelsAdapterInstance = &modelLabelsAdapter{daoImpl: dao}

	return nil
}

// ModelAdapter returns the instance of modelAdapter.
func ModelAdapter() *modelAdapter {
	return modelAdapterInstance
}

// ModelLabelsAdapter returns the instance of modelLabelsAdapter.
func ModelLabelsAdapter() *modelLabelsAdapter {
	return modelLabelsAdapterInstance
}
