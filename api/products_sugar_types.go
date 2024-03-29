package api

import (
	"github.com/applenperry-go/db"
	"github.com/applenperry-go/db/orm"
	"github.com/applenperry-go/model"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
)

func GetProductsSugarTypeByURL(c *gin.Context) {
	url := c.Param("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url param required"})
		return
	}
	var pst model.ProductsSugarType
	if err := db.DB.Where("url = ?", url).First(&pst).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pst)
}

func GetProductsSugarTypes(c *gin.Context) {
	var pst []model.ProductsSugarType
	if err := orm.GetList(db.DB.Preload("IconFile"), &pst, orm.Filters{
		Search:     c.Query("search"),
		SortColumn: c.Query("sort"),
		SortOrder:  c.Query("order"),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pst)
}

func GetProductsSugarType(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}
	var pst model.ProductsSugarType
	if err := orm.GetFirst(db.DB.Preload("IconFile"), &pst, id); err != nil {
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

	if err := db.DB.Updates(&pst).Error; err != nil {
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

	if err := db.DB.Delete(model.ProductsSugarType{ID: id}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "status": "deleted"})
}

func GetPossibleToDeleteProductsSugarType(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}

	var products []model.Product
	if err := db.DB.Where("sugar_type = ?", id).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(products) > 0 {
		deleteConflicts := make(map[string]interface{})
		deleteConflicts["products"] = products

		c.JSON(http.StatusOK, gin.H{
			"id":              id,
			"status":          "not_deletable",
			"deleteConflicts": deleteConflicts,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "status": "deletable"})
}
