package route

import (
	"github.com/applenperry-go/api"
	"github.com/applenperry-go/config"
	"github.com/gin-gonic/gin"
)

func Init(configuration config.Configuration) *gin.Engine {
	if configuration.PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.GET("/", api.Home)

	categories := r.Group("/categories")
	{
		categories.GET("/", api.GetCategories)
		categories.GET("/:id", api.GetCategory)
		categories.POST("/", api.CreateCategory)
		categories.PUT("/", api.UpdateCategory)
		categories.DELETE("/:id", api.DeleteCategory)
	}

	return r
}
