package spacerepositoryadapter

import "gorm.io/gorm"

var (
	spaceAdapterInstance       *spaceAdapter
	spaceLabelsAdapterInstance *spaceLabelsAdapter
)

func Init(db *gorm.DB, tables *Tables) error {
	// must set spaceTableName before migrating
	spaceTableName = tables.Space

	if err := db.AutoMigrate(&spaceDO{}); err != nil {
		return err
	}

	dbInstance = db

	dao := daoImpl{table: tables.Space}

	spaceAdapterInstance = &spaceAdapter{dao}
	spaceLabelsAdapterInstance = &spaceLabelsAdapter{dao}

	return nil
}

func SpaceAdapter() *spaceAdapter {
	return spaceAdapterInstance
}

func SpaceLabelsAdapter() *spaceLabelsAdapter {
	return spaceLabelsAdapterInstance
}
