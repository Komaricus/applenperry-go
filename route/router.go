package route

import (
	"github.com/applenperry-go/api"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func Init(db *gorm.DB) *echo.Echo {
	e := echo.New()

	e.GET("/", api.Home)
	e.GET("/categories", api.GetCategories(db))
	return e
}
