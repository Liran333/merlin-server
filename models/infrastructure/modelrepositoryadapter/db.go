package modelrepositoryadapter

import "gorm.io/gorm"

var (
	modelAdapterInstance       *modelAdapter
	modelLabelsAdapterInstance *modelLabelsAdapter
)

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

func ModelAdapter() *modelAdapter {
	return modelAdapterInstance
}

func ModelLabelsAdapter() *modelLabelsAdapter {
	return modelLabelsAdapterInstance
}
