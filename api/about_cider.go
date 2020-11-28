package api

import (
	"github.com/applenperry-go/db"
	"github.com/applenperry-go/db/orm"
	"github.com/applenperry-go/model"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
	"strings"
)

func GetWords(c *gin.Context) {
	var words []model.Word
	if err := db.DB.Find(&words).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, words)
}

func GetWordByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}
	var aboutCider model.AboutCider
	if err := orm.GetFirst(db.DB, &aboutCider, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, aboutCider)
}

func GetAboutCiderList(c *gin.Context) {
	var aboutCiderList []model.AboutCider
	if err := orm.GetList(db.DB, &aboutCiderList, orm.Filters{
		Search:     c.Query("search"),
		SortColumn: c.Query("sort"),
		SortOrder:  c.Query("order"),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, aboutCiderList)
}

func GetAboutCider(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}
	var aboutCider model.AboutCider
	if err := orm.GetFirst(db.DB, &aboutCider, id); err != nil {
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

	if err := db.DB.Where("cider_id = ?", aboutCider.ID).Delete(model.CiderAndFile{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := findImagesInDescription(aboutCider.ID, aboutCider.Description); err != nil {
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

	if err := db.DB.Updates(&aboutCider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Where("cider_id = ?", aboutCider.ID).Delete(model.CiderAndFile{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := findImagesInDescription(aboutCider.ID, aboutCider.Description); err != nil {
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

	if err := db.DB.Where("cider_id = ?", id).Delete(model.CiderAndFile{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Delete(model.AboutCider{ID: id}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "status": "deleted"})
}

func findImagesInDescription(ciderID, description string) error {
	const (
		start = "<img src=\"/images/"
		end   = "\">"
	)

	s := strings.Index(description, start)
	if s == -1 {
		return nil
	}
	s += len(start)

	e := strings.Index(description[s:], end)
	if e == -1 {
		return nil
	}
	e += s

	path := description[s:e]
	var file model.File
	if err := db.DB.Where("path = ?", path).First(&file).Error; err != nil {
		return err
	}

	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	if err := db.DB.Create(model.CiderAndFile{
		ID:      id.String(),
		CiderID: ciderID,
		FileID:  file.ID,
	}).Error; err != nil {
		return err
	}

	if err := findImagesInDescription(ciderID, description[e+len(end):]); err != nil {
		return err
	}

	return nil
}
