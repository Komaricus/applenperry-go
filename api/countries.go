package api

import (
	"github.com/applenperry-go/db"
	"github.com/applenperry-go/db/orm"
	"github.com/applenperry-go/model"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
)

func GetCountryByURL(c *gin.Context) {
	url := c.Param("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url param required"})
		return
	}
	var country model.Country
	if err := db.DB.Where("url = ?", url).First(&country).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, country)
}

func GetCountries(c *gin.Context) {
	var countries []model.Country
	if err := orm.GetList(db.DB.Preload("File").Preload("IconFile"), &countries, orm.Filters{
		Search:     c.Query("search"),
		SortColumn: c.Query("sort"),
		SortOrder:  c.Query("order"),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, countries)
}

func GetCountry(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}
	var country model.Country
	if err := orm.GetFirst(db.DB.Preload("File").Preload("IconFile"), &country, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, country)
}

func CreateCountry(c *gin.Context) {
	var country model.Country
	if err := c.Bind(&country); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	country.ID = id.String()

	if err := db.DB.Create(&country).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, country)
}

func UpdateCountry(c *gin.Context) {
	var country model.Country
	if err := c.Bind(&country); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Updates(&country).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, country)
}

func DeleteCountry(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}

	if err := db.DB.Delete(model.Country{ID: id}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "status": "deleted"})
}

func GetPossibleToDeleteCountry(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}

	var vendors []model.Vendor
	if err := db.DB.Where("country_id = ?", id).Find(&vendors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(vendors) > 0 {
		deleteConflicts := make(map[string]interface{})
		deleteConflicts["vendors"] = vendors
		c.JSON(http.StatusOK, gin.H{
			"id":              id,
			"status":          "not_deletable",
			"deleteConflicts": deleteConflicts,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "status": "deletable"})
}
