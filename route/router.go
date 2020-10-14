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
		appleApi.POST("/login", authMiddleware.LoginHandler)
		appleApi.GET("/refresh_token", authMiddleware.RefreshHandler)

		open := appleApi.Group("/open")
		{
			open.GET("/", api.Home)
			open.GET("slides", api.GetSlides)
		}

		categories := appleApi.Group("/categories")
		{
			categories.Use(authMiddleware.MiddlewareFunc())
			{
				categories.GET("/", api.GetCategories)
				categories.GET("/:id", api.GetCategory)
				categories.POST("/", api.CreateCategory)
				categories.PUT("/", api.UpdateCategory)
				categories.DELETE("/:id", api.DeleteCategory)
			}
		}

		aboutCider := appleApi.Group("/about-cider")
		{
			aboutCider.Use(authMiddleware.MiddlewareFunc())
			{
				aboutCider.GET("/", api.GetAboutCiderList)
				aboutCider.GET("/:id", api.GetAboutCider)
				aboutCider.POST("/", api.CreateAboutCider)
				aboutCider.PUT("/", api.UpdateAboutCider)
				aboutCider.DELETE("/:id", api.DeleteAboutCider)
			}
		}

		countries := appleApi.Group("/countries")
		{
			countries.Use(authMiddleware.MiddlewareFunc())
			{
				countries.GET("/", api.GetCountries)
				countries.GET("/:id", api.GetCountry)
				countries.POST("/", api.CreateCountry)
				countries.PUT("/", api.UpdateCountry)
				countries.DELETE("/:id", api.DeleteCountry)
			}
		}

		homeSlider := appleApi.Group("/home-slider")
		{
			homeSlider.Use(authMiddleware.MiddlewareFunc())
			{
				homeSlider.GET("/", api.GetHomeSliderItems)
				homeSlider.GET("/:id", api.GetHomeSliderItem)
				homeSlider.POST("/", api.CreateHomeSliderItem)
				homeSlider.PUT("/", api.UpdateHomeSliderItem)
				homeSlider.DELETE("/:id", api.DeleteHomeSliderItem)
			}
		}

		vendors := appleApi.Group("/vendors")
		{
			vendors.Use(authMiddleware.MiddlewareFunc())
			{
				vendors.GET("/", api.GetVendors)
				vendors.GET("/:id", api.GetVendor)
				vendors.POST("/", api.CreateVendor)
				vendors.PUT("/", api.UpdateVendor)
				vendors.DELETE("/:id", api.DeleteVendor)
			}
		}

		newsSections := appleApi.Group("/news-sections")
		{
			newsSections.Use(authMiddleware.MiddlewareFunc())
			{
				newsSections.GET("/", api.GetNewsSections)
				newsSections.GET("/:id", api.GetNewsSection)
				newsSections.POST("/", api.CreateNewsSection)
				newsSections.PUT("/", api.UpdateNewsSection)
				newsSections.DELETE("/:id", api.DeleteNewsSection)
			}
		}

		news := appleApi.Group("/news")
		{
			news.Use(authMiddleware.MiddlewareFunc())
			{
				news.GET("/", api.GetNews)
				news.GET("/:id", api.GetOneNews)
				news.POST("/", api.CreateNews)
				news.PUT("/", api.UpdateNews)
				news.DELETE("/:id", api.DeleteNews)
			}
		}

		files := appleApi.Group("/files")
		files.Use(authMiddleware.MiddlewareFunc())
		{
			files.POST("/upload", api.UploadFiles)
			files.GET("/", api.GetFiles)
			files.GET("/deletable/:id", api.GetPossibleToDeleteFile)
			files.DELETE("/", api.DeleteFile)
		}

		admins := appleApi.Group("/admins")
		admins.Use(authMiddleware.MiddlewareFunc())
		{
			admins.POST("/", api.CreateAdmin)
			admins.GET("/heartbeat", api.CheckHeartbeat)
		}
	}

	return r
}
