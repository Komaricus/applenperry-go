package api

import (
	"github.com/applenperry-go/db"
	"github.com/applenperry-go/model"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
	"strconv"
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

		var product model.Product
		if err := db.DB.Where("id = ?", p.ProductID).First(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		count := product.Amount - p.ProductCount
		if count < 0 {
			count = 0
		}
		if err := db.DB.Where("id = ?", p.ProductID).Update("amount", count).Error; err != nil {
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

func UpdateOrderStatus(c *gin.Context) {
	var update model.UpdateOrder
	if err := c.Bind(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Updates(&update).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, update)
}

func DeleteOrder(c *gin.Context) {
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

	for _, p := range order.Products {
		if err := deleteProductFromOrder(order.ID, p.ProductID, int(p.ProductCount)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := db.DB.Delete(model.DeleteOrder{ID: id}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "status": "deleted"})
}

func DeleteProductFromOrder(c *gin.Context) {
	productID := c.Query("productId")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "productId query param required"})
		return
	}

	orderID := c.Query("orderId")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "orderId query param required"})
		return
	}

	productAmount, err := strconv.Atoi(c.Query("productAmount"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := deleteProductFromOrder(orderID, productID, productAmount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"productId": productID, "orderId": orderID, "status": "deleted"})
}

func deleteProductFromOrder(orderID, productID string, productAmount int) error {
	if err := db.DB.Where("product_id = ?", productID).Where("order_id = ?", orderID).Delete(model.OrderAndProduct{}).Error; err != nil {
		return err
	}

	if productAmount > 0 {
		var product model.Product
		if err := db.DB.Where("id = ?", productID).First(&product).Error; err != nil {
			return err
		}

		amount := product.Amount + uint(productAmount)
		if err := db.DB.Model(model.Product{}).Where("id = ?", productID).Update("amount", amount).Error; err != nil {
			return err
		}
	}

	return nil
}
