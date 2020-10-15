package api

import (
	"github.com/applenperry-go/db"
	"github.com/applenperry-go/model"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
)

func GetProductsSugarTypes(c *gin.Context) {
	var pst []model.ProductsSugarType
	if err := db.DB.Where("is_deleted = false").Find(&pst).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pst)
}

func GetProductsSugarType(c *gin.Context) {
	var pst model.ProductsSugarType
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}

	if err := db.DB.Where("id = ?", id).Where("is_deleted = false").First(&pst).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pst)
}

func CreateProductsSugarType(c *gin.Context) {
	var pst model.ProductsSugarType
	if err := c.Bind(&pst); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	pst.ID = id.String()

	if err := db.DB.Create(&pst).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pst)
}

func UpdateProductsSugarType(c *gin.Context) {
	var pst model.ProductsSugarType
	if err := c.Bind(&pst); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Updates(model.ProductsSugarType{
		ID:   pst.ID,
		Name: pst.Name,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pst)
}

func DeleteProductsSugarType(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}

	if err := db.DB.Model(model.ProductsSugarType{ID: id}).Update("is_deleted", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "status": "deleted"})
}
