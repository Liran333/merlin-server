package postgresql

import (
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(cfg *Config) error {
	dbInstance, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN: cfg.dsn(),
			// disables implicit prepared statement usage
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{},
	)
	if err != nil {
		return err
	}

	sqlDb, err := dbInstance.DB()
	if err != nil {
		return err
	}

	sqlDb.SetConnMaxLifetime(cfg.getLifeDuration())
	sqlDb.SetMaxOpenConns(cfg.MaxConn)
	sqlDb.SetMaxIdleConns(cfg.MaxIdle)

	db = dbInstance

	return nil
}

func DB() *gorm.DB {
	return db
}

func AutoMigrate(table interface{}) error {
	// pointer non-nil check
	if db == nil {
		err := errors.New("empty pointer of *gorm.DB")
		logrus.Error(err.Error())

		return err
	}

	return db.AutoMigrate(table)
}
