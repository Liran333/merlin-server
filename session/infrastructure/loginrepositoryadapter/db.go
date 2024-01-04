package loginrepositoryadapter

import "github.com/openmerlin/merlin-server/common/infrastructure/postgresql"

var loginAdapterInstance *loginAdapter

func Init(tables *Tables) error {
	// must set loginTableName before migrating
	loginTableName = tables.Login

	if err := postgresql.AutoMigrate(&loginDO{}); err != nil {
		return err
	}

	dao := postgresql.DAO(tables.Login)

	loginAdapterInstance = &loginAdapter{dao}

	return nil
}

func LoginAdapter() *loginAdapter {
	return loginAdapterInstance
}
