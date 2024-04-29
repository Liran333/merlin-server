/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package repositoryadapter

import "gorm.io/gorm"

var (
	computilityAdapterInstance              *computilityOrgAdapter
	computilityDetailAdapterInstance        *computilityDetailAdapter
	computilityAccountAdapterInstance       *computilityAccountAdapter
	computilityAccountRecordAdapterInstance *computilityAccountRecordAdapter
)

// Init initializes the database and sets up the necessary adapters.
func Init(db *gorm.DB, tables *Tables) error {
	// must set TableName before migrating
	computilityOrgTableName = tables.ComputilityOrg
	computilityDetailTableName = tables.ComputilityDetail
	computilityAccountTableName = tables.ComputilityAccount
	computilityAccountRecordTableName = tables.ComputilityAccountRecord

	if err := db.AutoMigrate(&computilityOrgDO{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&computilityDetailDO{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&computilityAccountDO{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&computilityAccountRecordDO{}); err != nil {
		return err
	}

	dbInstance = db

	computilityDao := daoImpl{table: computilityOrgTableName}
	computilityDetailDao := daoImpl{table: computilityDetailTableName}
	computilityAccountDao := daoImpl{table: computilityAccountTableName}
	computilityAccountRecordDao := daoImpl{table: computilityAccountRecordTableName}

	computilityAdapterInstance = &computilityOrgAdapter{
		daoImpl: computilityDao,
	}
	computilityDetailAdapterInstance = &computilityDetailAdapter{
		daoImpl: computilityDetailDao,
	}
	computilityAccountAdapterInstance = &computilityAccountAdapter{
		daoImpl: computilityAccountDao,
	}
	computilityAccountRecordAdapterInstance = &computilityAccountRecordAdapter{
		daoImpl: computilityAccountRecordDao,
	}

	return nil
}

func ComputilityOrgAdapter() *computilityOrgAdapter {
	return computilityAdapterInstance
}

func ComputilityDetailAdapter() *computilityDetailAdapter {
	return computilityDetailAdapterInstance
}

func ComputilityAccountAdapter() *computilityAccountAdapter {
	return computilityAccountAdapterInstance
}

func ComputilityAccountRecordAdapter() *computilityAccountRecordAdapter {
	return computilityAccountRecordAdapterInstance
}
