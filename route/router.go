package route

import (
	"github.com/applenperry-go/api"
	"github.com/applenperry-go/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init(configuration config.Configuration) *gin.Engine {
	if configuration.PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	conf := cors.DefaultConfig()
	conf.AllowOrigins = []string{"https://applenperry.ru", "https://www.applenperry.ru"}

	r.Use(cors.New(conf))

	appleApi := r.Group("/apple-api")
	{
		appleApi.GET("/", api.Home)
		categories := appleApi.Group("/categories")
		{
			categories.GET("/", api.GetCategories)
			categories.GET("/:id", api.GetCategory)
			categories.POST("/", api.CreateCategory)
			categories.PUT("/", api.UpdateCategory)
			categories.DELETE("/:id", api.DeleteCategory)
		}
	}

	return r
}
