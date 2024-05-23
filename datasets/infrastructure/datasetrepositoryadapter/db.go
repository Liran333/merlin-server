/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package datasetrepositoryadapter provides an adapter for the datasets repository
package datasetrepositoryadapter

import "gorm.io/gorm"

var (
	datasetAdapterInstance       *datasetAdapter
	datasetLabelsAdapterInstance *datasetLabelsAdapter
)

// Init initializes the dataset module by performing necessary setup and migrations.
func Init(db *gorm.DB, tables *Tables) error {
	// must set datasetTableName before migrating
	datasetTableName = tables.Datasets

	if err := db.AutoMigrate(&datasetDO{}); err != nil {
		return err
	}

	dbInstance = db

	dao := daoImpl{table: tables.Datasets}

	datasetAdapterInstance = &datasetAdapter{daoImpl: dao}
	datasetLabelsAdapterInstance = &datasetLabelsAdapter{daoImpl: dao}

	return nil
}

// DatasetAdapter returns the instance of datasetAdapter.
func DatasetAdapter() *datasetAdapter {
	return datasetAdapterInstance
}

// DatasetLabelsAdapter returns the instance of datasetLabelsAdapter.
func DatasetLabelsAdapter() *datasetLabelsAdapter {
	return datasetLabelsAdapterInstance
}
