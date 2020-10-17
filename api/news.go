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

func GetNews(c *gin.Context) {
	var news []model.News
	q := db.DB.Preload("File").Preload("Section")
	if err := orm.GetList(q, &news, orm.Filters{
		Search:     c.Query("search"),
		SortColumn: c.Query("sort"),
		SortOrder:  c.Query("order"),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, news)
}

func GetOneNews(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}
	var news model.News
	q := db.DB.Preload("File").Preload("Section")
	if err := orm.GetFirst(q, &news, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, news)
}

func CreateNews(c *gin.Context) {
	var news model.News
	if err := c.Bind(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	news.ID = id.String()

	if err := db.DB.Create(&news).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Where("news_id = ?", news.ID).Delete(model.NewsAndFiles{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := findImages(news.ID, news.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, news)
}

func UpdateNews(c *gin.Context) {
	var news model.News
	if err := c.Bind(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Updates(&news).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Where("news_id = ?", news.ID).Delete(model.NewsAndFiles{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := findImages(news.ID, news.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, news)
}

func DeleteNews(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}

	if err := db.DB.Where("news_id = ?", id).Delete(model.NewsAndFiles{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Delete(model.News{ID: id}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "status": "deleted"})
}

func findImages(newsId, content string) error {
	const (
		start = "<img src=\"/images/"
		end   = "\">"
	)

	s := strings.Index(content, start)
	if s == -1 {
		return nil
	}
	s += len(start)

	e := strings.Index(content, end)
	if e == -1 {
		return nil
	}

	path := content[s:e]
	var file model.File
	if err := db.DB.Where("path = ?", path).First(&file).Error; err != nil {
		return err
	}

	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	if err := db.DB.Create(model.NewsAndFiles{
		ID:     id.String(),
		NewsID: newsId,
		FileID: file.ID,
	}).Error; err != nil {
		return err
	}

	if err := findImages(newsId, content[e+len(end):]); err != nil {
		return err
	}

	return nil
}
