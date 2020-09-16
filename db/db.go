package db

import (
	"fmt"
	"github.com/applenperry-go/config"
	"github.com/applenperry-go/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	configuration := config.GetConfig()
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=posix/Europe/Moscow", configuration.DB_HOST, configuration.DB_USERNAME, configuration.DB_PASSWORD, configuration.DB_NAME, configuration.DB_PORT)
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic("DB Connection Error")
	}

	initialMigration(db)

	return db
}

func initialMigration(db *gorm.DB) {
	db.AutoMigrate(&model.Category{})
}
