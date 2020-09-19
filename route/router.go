package route

import (
	"github.com/applenperry-go/api"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/", api.Home)

	categories := r.Group("/categories")
	{
		categories.GET("/", api.GetCategories)
		categories.GET("/:id", api.GetCategory)
		categories.POST("/", api.CreateCategory)
	}

	return r
}
