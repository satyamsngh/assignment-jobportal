package database

import (
	"fmt"
	"job-portal-api/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open(config config.DataConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.Host, config.UserName, config.Password, config.DBName, config.Port, config.SSLMode, config.Time)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
