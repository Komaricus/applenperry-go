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

	origins := []string{"https://applenperry.ru", "https://www.applenperry.ru"}
	if !configuration.PRODUCTION {
		origins = append(origins, "http://localhost:8080")
	}
	conf.AllowOrigins = origins
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
			open.GET("/slides", api.GetSlides)
			open.GET("/shop-slides", api.GetShopSlides)
			open.GET("/products", api.GetProductsWithPaginate)
			open.GET("/categories", api.GetCategoriesWithChild)
			open.GET("/categories/:url", api.GetCategoryByURL)
			open.GET("/countries/:url", api.GetCountryByURL)
			open.GET("/products-types/:url", api.GetProductsTypeByURL)
			open.GET("/products-sugar-types/:url", api.GetProductsSugarTypeByURL)
			open.GET("/products/:url", api.GetProductByURL)
			open.GET("/vendors", api.GetVendorsList)
			open.GET("/vendors/:url", api.GetVendorByURL)
			open.POST("/order", api.CreateOrder)
			open.GET("/docs/:url", api.GetDocumentByURL)
			open.GET("/docs", api.GetOpenDocs)
			open.GET("/pages/:url", api.GetPageByURL)
			open.GET("/words", api.GetWords)
			open.GET("/words/:id", api.GetWordByID)
		}

		orders := appleApi.Group("/orders")
		{
			orders.GET("/", api.GetOrders)
			orders.GET("/:id", api.GetOrder)
			orders.DELETE("/product", api.DeleteProductFromOrder)
			orders.DELETE("/order/:id", api.DeleteOrder)
			orders.PUT("/", api.UpdateOrderStatus)
		}

		docs := appleApi.Group("/docs")
		{
			docs.GET("/", api.GetDocsList)
			docs.GET("/:id", api.GetDocument)
			docs.POST("/", api.CreateDocument)
			docs.PUT("/", api.UpdateDocument)
			docs.DELETE("/:id", api.DeleteDocument)
		}

		pages := appleApi.Group("/pages")
		{
			pages.GET("/", api.GetPages)
			pages.GET("/:id", api.GetPage)
			pages.POST("/", api.CreatePage)
			pages.PUT("/", api.UpdatePage)
			pages.DELETE("/:id", api.DeletePage)
		}

		categories := appleApi.Group("/categories")
		{
			categories.Use(authMiddleware.MiddlewareFunc())
			{
				categories.GET("/", api.GetCategories)
				categories.GET("/:id", api.GetCategory)
				categories.GET("/:id/deletable", api.GetPossibleToDeleteCategory)
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
				countries.GET("/:id/deletable", api.GetPossibleToDeleteCountry)
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

		shopSlider := appleApi.Group("/shop-slider")
		{
			shopSlider.Use(authMiddleware.MiddlewareFunc())
			{
				shopSlider.GET("/", api.GetShopSliderItems)
				shopSlider.GET("/:id", api.GetShopSliderItem)
				shopSlider.POST("/", api.CreateShopSliderItem)
				shopSlider.PUT("/", api.UpdateShopSliderItem)
				shopSlider.DELETE("/:id", api.DeleteShopSliderItem)
			}
		}

		vendors := appleApi.Group("/vendors")
		{
			vendors.Use(authMiddleware.MiddlewareFunc())
			{
				vendors.GET("/", api.GetVendors)
				vendors.GET("/:id", api.GetVendor)
				vendors.GET("/:id/deletable", api.GetPossibleToDeleteVendor)
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
				newsSections.GET("/:id/deletable", api.GetPossibleToDeleteNewsSection)
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

		productsTypes := appleApi.Group("/products-types")
		{
			productsTypes.Use(authMiddleware.MiddlewareFunc())
			{
				productsTypes.GET("/", api.GetProductsTypes)
				productsTypes.GET("/:id", api.GetProductsType)
				productsTypes.GET("/:id/deletable", api.GetPossibleToDeleteProductsType)
				productsTypes.POST("/", api.CreateProductsType)
				productsTypes.PUT("/", api.UpdateProductsType)
				productsTypes.DELETE("/:id", api.DeleteProductsType)
			}
		}

		productsSugarTypes := appleApi.Group("/products-sugar-types")
		{
			productsSugarTypes.Use(authMiddleware.MiddlewareFunc())
			{
				productsSugarTypes.GET("/", api.GetProductsSugarTypes)
				productsSugarTypes.GET("/:id", api.GetProductsSugarType)
				productsSugarTypes.GET("/:id/deletable", api.GetPossibleToDeleteProductsSugarType)
				productsSugarTypes.POST("/", api.CreateProductsSugarType)
				productsSugarTypes.PUT("/", api.UpdateProductsSugarType)
				productsSugarTypes.DELETE("/:id", api.DeleteProductsSugarType)
			}
		}

		products := appleApi.Group("/products")
		{
			products.Use(authMiddleware.MiddlewareFunc())
			{
				products.GET("/", api.GetProducts)
				products.GET("/:id", api.GetProduct)
				products.POST("/", api.CreateProduct)
				products.PUT("/", api.UpdateProduct)
				products.DELETE("/:id", api.DeleteProduct)
				products.GET("/:id/deletable", api.GetPossibleToDeleteProduct)
			}
		}

		files := appleApi.Group("/files")
		{
			files.Use(authMiddleware.MiddlewareFunc())
			{
				files.POST("/upload", api.UploadFiles)
				files.GET("/", api.GetFiles)
				files.GET("/deletable/:id", api.GetPossibleToDeleteFile)
				files.DELETE("/", api.DeleteFile)
				files.GET("/download/:id", api.DownloadFile)
			}
		}

		admins := appleApi.Group("/admins")
		{
			admins.Use(authMiddleware.MiddlewareFunc())
			{
				admins.POST("/", api.CreateAdmin)
				admins.GET("/heartbeat", api.CheckHeartbeat)
			}
		}
	}

	return r
}
