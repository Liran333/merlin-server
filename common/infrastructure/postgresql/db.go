package postgresql

import (
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var (
	db         *gorm.DB
	errorCodes errorCode
)

var serverLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	logger.Config{
		LogLevel:                  logger.Warn, // Log level
		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
		ParameterizedQueries:      true,        // Don't include params in the SQL log
		Colorful:                  false,       // Disable color
	},
)

func Init(cfg *Config, removeCfg bool) error {
	dbInstance, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN: cfg.dsn(),
			// disables implicit prepared statement usage
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{
			Logger: serverLogger,
		},
	)
	if err != nil {
		return err
	}

	if removeCfg && cfg.Dbcert != "" {
		if err := os.Remove(cfg.Dbcert); err != nil {
			return err
		}
	}

	sqlDb, err := dbInstance.DB()
	if err != nil {
		return err
	}

	sqlDb.SetConnMaxLifetime(cfg.getLifeDuration())
	sqlDb.SetMaxOpenConns(cfg.MaxConn)
	sqlDb.SetMaxIdleConns(cfg.MaxIdle)

	db = dbInstance

	errorCodes = cfg.Code

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
