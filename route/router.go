package route

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/applenperry-go/api"
	"github.com/applenperry-go/api/middleware"
	"github.com/applenperry-go/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

func Init(configuration config.Configuration) *gin.Engine {
	if configuration.PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	conf := cors.DefaultConfig()
	conf.AllowOrigins = []string{"https://applenperry.ru", "https://www.applenperry.ru"}
	conf.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	r.Use(cors.New(conf))

	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	authMiddleware, err := middleware.GetAuthMiddleware()
	if err != nil {
		log.Fatal(err.Error())
	}

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	appleApi := r.Group("/apple-api")
	{
		appleApi.GET("/", api.Home)
		appleApi.POST("/login", authMiddleware.LoginHandler)
		appleApi.GET("/refresh_token", authMiddleware.RefreshHandler)

		categories := appleApi.Group("/categories")
		{
			categories.GET("/", api.GetCategories)
			categories.GET("/:id", api.GetCategory)
			categories.Use(authMiddleware.MiddlewareFunc())
			{
				categories.POST("/", api.CreateCategory)
				categories.PUT("/", api.UpdateCategory)
				categories.DELETE("/:id", api.DeleteCategory)
			}
		}

		aboutCider := appleApi.Group("/about-cider")
		{
			aboutCider.GET("/", api.GetAboutCiderList)
			aboutCider.GET("/:id", api.GetAboutCider)
			aboutCider.Use(authMiddleware.MiddlewareFunc())
			{
				aboutCider.POST("/", api.CreateAboutCider)
				aboutCider.PUT("/", api.UpdateAboutCider)
				aboutCider.DELETE("/:id", api.DeleteAboutCider)
			}
		}

		files := appleApi.Group("/files")
		files.Use(authMiddleware.MiddlewareFunc())
		{
			files.POST("/upload", api.UploadFiles)
		}

		admins := appleApi.Group("/admins")
		admins.Use(authMiddleware.MiddlewareFunc())
		{
			admins.POST("/", api.CreateAdmin)
		}
	}

	return r
}
