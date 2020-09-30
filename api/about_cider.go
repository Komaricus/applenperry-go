package api

import (
	"github.com/applenperry-go/db"
	"github.com/applenperry-go/model"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
)

func GetAboutCiderList(c *gin.Context) {
	var aboutCiderList []model.AboutCider
	if err := db.DB.Where("is_deleted = false").Find(&aboutCiderList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, aboutCiderList)
}

func GetAboutCider(c *gin.Context) {
	var aboutCider model.AboutCider
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}

	if err := db.DB.Where("id = ?", id).Where("is_deleted = false").First(&aboutCider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, aboutCider)
}

func CreateAboutCider(c *gin.Context) {
	var aboutCider model.AboutCider
	if err := c.Bind(&aboutCider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	aboutCider.ID = id.String()

	if err := db.DB.Create(&aboutCider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, aboutCider)
}

func UpdateAboutCider(c *gin.Context) {
	var aboutCider model.AboutCider
	if err := c.Bind(&aboutCider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Updates(model.AboutCider{
		ID:          aboutCider.ID,
		Name:        aboutCider.Name,
		Description: aboutCider.Description,
		Size:        aboutCider.Size,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, aboutCider)
}

func DeleteAboutCider(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}

	if err := db.DB.Model(model.AboutCider{ID: id}).Update("is_deleted", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "status": "deleted"})
}
