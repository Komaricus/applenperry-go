package api

import (
	"github.com/applenperry-go/db"
	"github.com/applenperry-go/db/orm"
	"github.com/applenperry-go/model"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
)

func GetVendorByURL(c *gin.Context) {
	url := c.Param("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url param required"})
		return
	}
	var v model.Vendor
	if err := db.DB.Where("url = ?", url).First(&v).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, v)
}

func GetVendorsList(c *gin.Context) {
	var vendors []model.Vendor
	q := db.DB.Preload("File").Preload("Country")
	if err := q.Order("name").Find(&vendors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vendors)
}

func GetVendors(c *gin.Context) {
	var vendors []model.Vendor
	q := db.DB.Preload("File").Preload("Country")
	if err := orm.GetList(q, &vendors, orm.Filters{
		Search:     c.Query("search"),
		SortColumn: c.Query("sort"),
		SortOrder:  c.Query("order"),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vendors)
}

func GetVendor(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}
	var vendor model.Vendor
	q := db.DB.Preload("File").Preload("Country")
	if err := orm.GetFirst(q, &vendor, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vendor)
}

func CreateVendor(c *gin.Context) {
	var vendor model.Vendor
	if err := c.Bind(&vendor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	vendor.ID = id.String()

	if err := db.DB.Create(&vendor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, vendor)
}

func UpdateVendor(c *gin.Context) {
	var vendor model.Vendor
	if err := c.Bind(&vendor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Updates(&vendor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, vendor)
}

func DeleteVendor(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}

	if err := db.DB.Delete(model.Vendor{ID: id}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "status": "deleted"})
}

func GetPossibleToDeleteVendor(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}

	var products []model.Product
	if err := db.DB.Where("vendor_id = ?", id).Find(&products).Error; err != nil {
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
