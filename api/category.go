package api

import (
	"github.com/applenperry-go/model"
	"github.com/labstack/echo"
	"gorm.io/gorm"
	"net/http"
)

func GetCategories(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		var categories []model.Category
		db.Find(&categories)
		return c.JSON(http.StatusOK, categories)
	}
}
