/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package modelrepositoryadapter provides an adapter for the model repository
package modelrepositoryadapter

import "gorm.io/gorm"

var (
	modelAdapterInstance       *modelAdapter
	modelLabelsAdapterInstance *modelLabelsAdapter
	modelDeployAdapterInstance *modelDeployAdapter
)

// Init initializes the model module by performing necessary setup and migrations.
func Init(db *gorm.DB, tables *Tables) error {
	// must set modelTableName before migrating
	modelTableName = tables.Model
	modelDeployTableName = tables.ModelDeploy

	if err := db.AutoMigrate(&modelDO{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&modelDeployDO{}); err != nil {
		return err
	}

	dbInstance = db

	dao := daoImpl{table: tables.Model}
	daoDeploy := daoImpl{table: tables.ModelDeploy}

	modelAdapterInstance = &modelAdapter{daoImpl: dao}
	modelLabelsAdapterInstance = &modelLabelsAdapter{daoImpl: dao}
	modelDeployAdapterInstance = &modelDeployAdapter{daoImpl: daoDeploy}

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

func ModelDeployAdapter() *modelDeployAdapter {
	return modelDeployAdapterInstance
}
