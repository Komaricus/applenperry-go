package api

import (
	"github.com/applenperry-go/db"
	"github.com/applenperry-go/model"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
)

func GetSlides(c *gin.Context) {
	var slides []model.Slide
	q := db.DB.Preload("File").Where("dbo.home_slider.is_deleted = false").Order("priority")
	if err := q.Find(&slides).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, slides)
}

func GetHomeSliderItems(c *gin.Context) {
	var items []model.HomeSliderItem
	q := db.DB.Preload("File").Where("dbo.home_slider.is_deleted = false")
	if err := q.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func GetHomeSliderItem(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}

	var item model.HomeSliderItem

	q := db.DB.Preload("File").Where("dbo.home_slider.is_deleted = false").Where("id = ?", id)

	if err := q.First(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func CreateHomeSliderItem(c *gin.Context) {
	var item model.HomeSliderItem
	if err := c.Bind(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	item.ID = id.String()

	if err := db.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func UpdateHomeSliderItem(c *gin.Context) {
	var item model.HomeSliderItem
	if err := c.Bind(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Updates(model.HomeSliderItem{
		ID:       item.ID,
		Name:     item.Name,
		Priority: item.Priority,
		FileID:   item.FileID,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func DeleteHomeSliderItem(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}

	if err := db.DB.Model(model.HomeSliderItem{ID: id}).Update("is_deleted", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "status": "deleted"})
}
