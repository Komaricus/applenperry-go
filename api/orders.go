package api

import (
	"github.com/applenperry-go/db"
	"github.com/applenperry-go/model"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
)

func GetOrders(c *gin.Context) {
	var orders []model.GetOrder
	q := db.DB.Preload("Products").Preload("Products.Product").Preload("Products.Product.MainImage")
	if err := q.Order("code desc").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func GetOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}

	var order model.GetOrder
	q := db.DB.Preload("Products").Preload("Products.Product").Preload("Products.Product.MainImage")
	if err := q.Where("id = ?", id).Find(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func CreateOrder(c *gin.Context) {
	var order model.CreateOrder
	if err := c.Bind(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	order.ID = id.String()
	order.Status = "new"

	if err := db.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, p := range order.Products {
		if err := db.DB.Create(&model.OrderAndProduct{
			OrderID:      order.ID,
			ProductID:    p.ProductID,
			ProductCount: p.ProductCount,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	var resp model.GetOrder
	if err := db.DB.Where("id = ?", order.ID).Find(&resp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}
