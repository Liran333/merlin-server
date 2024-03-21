/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package spacerepositoryadapter

import "gorm.io/gorm"

var (
	spaceAdapterInstance         *spaceAdapter
	spaceLabelsAdapterInstance   *spaceLabelsAdapter
	spaceModelInstance           *modelSpaceRelationAdapter
	spaceVariableAdapterInstance *spaceVariableAdapter
	spaceSecretAdapterInstance   *spaceSecretAdapter
)

// Init initializes the database and sets up the necessary adapters.
func Init(db *gorm.DB, tables *Tables) error {
	// must set spaceTableName before migrating
	spaceTableName = tables.Space
	spaceModelRelationTableName = tables.SpaceModel
	spaceEnvSecretTableName = tables.SpaceEnvSecret

	if err := db.AutoMigrate(&spaceDO{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&modelSpaceRelationDO{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&spaceEnvSecretDO{}); err != nil {
		return err
	}

	dbInstance = db

	spaceDao := daoImpl{table: tables.Space}
	spaceModelDao := daoImpl{table: tables.SpaceModel}
	spaceEnvSecretDao := daoImpl{table: tables.SpaceEnvSecret}

	spaceAdapterInstance = &spaceAdapter{daoImpl: spaceDao}
	spaceLabelsAdapterInstance = &spaceLabelsAdapter{daoImpl: spaceDao}
	spaceModelInstance = &modelSpaceRelationAdapter{daoImpl: spaceModelDao}
	spaceVariableAdapterInstance = &spaceVariableAdapter{daoImpl: spaceEnvSecretDao}
	spaceSecretAdapterInstance = &spaceSecretAdapter{daoImpl: spaceEnvSecretDao}

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

// SpaceVariableAdapter returns the instance of the space variable adapter.
func SpaceVariableAdapter() *spaceVariableAdapter {
	return spaceVariableAdapterInstance
}

// SpaceSecretAdapter returns the instance of the space secret adapter.
func SpaceSecretAdapter() *spaceSecretAdapter {
	return spaceSecretAdapterInstance
}
