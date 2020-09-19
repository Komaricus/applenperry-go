package db

import (
	"fmt"
	"github.com/applenperry-go/config"
	"github.com/applenperry-go/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(configuration config.Configuration) error {
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=posix/Europe/Moscow", configuration.DB_HOST, configuration.DB_USERNAME, configuration.DB_PASSWORD, configuration.DB_NAME, configuration.DB_PORT)
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return err
	}
	initialMigration(db)
	DB = db
	return nil
}

func initialMigration(db *gorm.DB) {
	db.AutoMigrate(&model.Category{})
}
