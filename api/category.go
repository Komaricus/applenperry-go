package api

import (
	"github.com/applenperry-go/db"
	"github.com/applenperry-go/model"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
)

func GetCategories(c *gin.Context) {
	var categories []model.Category
	if err := db.DB.Where("is_deleted = false").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func GetCategory(c *gin.Context) {
	var category model.Category
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}

	if err := db.DB.Where("id = ?", id).Where("is_deleted = false").First(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

func CreateCategory(c *gin.Context) {
	var category model.Category
	if err := c.Bind(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	category.ID = id.String()

	if err := db.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func UpdateCategory(c *gin.Context) {
	var category model.Category
	if err := c.Bind(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if category.ParentID == nil {
		if err := db.DB.Model(category).UpdateColumn("parent_id", nil).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if err := db.DB.Updates(model.Category{
		ID:          category.ID,
		Name:        category.Name,
		Url:         category.Url,
		Description: category.Description,
		ParentID:    category.ParentID,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}

	if err := db.DB.Model(model.Category{ID: id}).Update("is_deleted", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := DeleteChildCategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "status": "deleted"})
}

func DeleteChildCategory(id string) error {
	var child []model.Category
	if err := db.DB.Where("is_deleted = false").Where("parent_id = ?", id).Find(&child).Error; err != nil {
		return err
	}

	for _, c := range child {
		if err := db.DB.Model(model.Category{ID: c.ID}).Update("is_deleted", true).Error; err != nil {
			return err
		}
		if err := DeleteChildCategory(c.ID); err != nil {
			return err
		}
	}

	return nil
}
