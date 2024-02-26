/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package loginrepositoryadapter

import (
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
)

var loginAdapterInstance *loginAdapter

// Init initializes the login adapter with the given tables.
func Init(tables *Tables) error {
	// must set loginTableName before migrating
	loginTableName = tables.Login

	if err := postgresql.AutoMigrate(&loginDO{}); err != nil {
		return err
	}

	dao := postgresql.DAO(tables.Login)

	loginAdapterInstance = &loginAdapter{dao: dao}

	return nil
}

// LoginAdapter returns the login adapter instance.
func LoginAdapter() *loginAdapter {
	return loginAdapterInstance
}
