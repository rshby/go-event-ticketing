package database

import (
	"fmt"

	"github.com/rshby/go-event-ticketing/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectPostgreSql connect to SQL
func ConnectPostgreSql() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s connect_timeout=%d",
		config.DbHost(), config.DbUser(), config.DbPassword(), config.DbName(), config.DbPort(), config.DbTimezone(), int(config.DbConnectionTimeout().Seconds()))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Info),
		TranslateError:         true,
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	dbSql, err := db.DB()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// config connection pool
	dbSql.SetMaxOpenConns(config.DbMaxOpenConns())
	dbSql.SetMaxIdleConns(config.DbMaxIdleConns())
	dbSql.SetConnMaxLifetime(config.DbConnMaxLifetime())
	dbSql.SetConnMaxIdleTime(config.DbConnMaxIdletime())

	logrus.Infof("success connect to database [%s]✅", dsn)

	return db, nil
}
