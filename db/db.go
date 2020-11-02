package db

import (
	"fmt"
	"github.com/applenperry-go/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(configuration config.Configuration) error {
	var connectionString string
	if configuration.PRODUCTION {
		connectionString = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=posix/Europe/Moscow", configuration.DB_HOST, configuration.DB_USERNAME, configuration.DB_PASSWORD, configuration.DB_NAME, configuration.DB_PORT)
	} else {
		connectionString = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", configuration.DB_HOST, configuration.DB_USERNAME, configuration.DB_PASSWORD, configuration.DB_NAME, configuration.DB_PORT)
	}

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	return nil
}
