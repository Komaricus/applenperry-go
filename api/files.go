package api

import (
	"fmt"
	"github.com/applenperry-go/db"
	"github.com/applenperry-go/model"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func GetFiles(c *gin.Context) {
	var files []model.File

	page := 1
	pageParam := c.Query("page")
	if pageParam != "" {
		page, _ = strconv.Atoi(pageParam)
	}

	perPage := -1
	perPageParam := c.Query("perPage")
	if perPageParam != "" {
		perPage, _ = strconv.Atoi(perPageParam)
	}

	search := c.Query("search")

	offset := (page - 1) * perPage

	q := db.DB.Offset(offset).Order("created_at desc")
	t := db.DB.Model(model.File{}).Group("id")

	if perPage != -1 {
		q.Limit(perPage)
	}

	if search != "" {
		search = "%" + search + "%"
		q.Where("original_name LIKE ?", search)
		t.Where("original_name LIKE ?", search)
	}

	if err := q.Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var count int64
	if err := t.Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.GetFilesResponse{
		Total: count,
		Files: files,
	})
}

func UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	files := form.File["files"]

	dbFiles := make([]model.File, 0)
	basePath := os.Getenv("IMAGES_PATH")

	for _, file := range files {
		id, err := uuid.NewV4()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		originalName := filepath.Base(file.Filename)
		ext := filepath.Ext(file.Filename)
		filename := id.String() + ext

		filePath, err := generateFilePath(basePath)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		path := basePath + filePath + filename

		if err := c.SaveUploadedFile(file, path); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		dbFile := model.File{
			ID:           id.String(),
			FileName:     filename,
			Path:         filePath + filename,
			OriginalName: originalName,
			Size:         file.Size,
		}

		if err := db.DB.Create(&dbFile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		dbFiles = append(dbFiles, dbFile)
	}
	c.JSON(http.StatusOK, dbFiles)
}

func GetPossibleToDeleteFile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}

	var countries []model.Country
	if err := db.DB.Where("flag = ?", id).Find(&countries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var slides []model.HomeSliderItem
	if err := db.DB.Where("file_id = ?", id).Find(&slides).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var vendors []model.Vendor
	if err := db.DB.Where("file_id = ?", id).Find(&vendors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var news []model.News
	if err := db.DB.Where("file_id = ?", id).Find(&news).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newsMap := make(map[string]bool, len(news))
	for _, n := range news {
		if ok := newsMap[n.ID]; !ok {
			newsMap[n.ID] = true
		}
	}

	var naf []model.NewsAndFiles
	if err := db.DB.Joins("News").Where("dbo.news_and_files.file_id = ?", id).Find(&naf).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, nf := range naf {
		if ok := newsMap[nf.NewsID]; !ok {
			newsMap[nf.NewsID] = true
			news = append(news, nf.News)
		}
	}

	var paf []model.ProductsAndFiles
	if err := db.DB.Where("file_id = ?", id).Find(&paf).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ids := make([]string, 0, len(paf))
	for _, n := range paf {
		ids = append(ids, n.ProductID)
	}

	var products []model.Product
	if err := db.DB.Where("id IN (?)", ids).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(countries) > 0 || len(slides) > 0 || len(vendors) > 0 || len(news) > 0 || len(products) > 0 {
		deleteConflicts := make(map[string]interface{})
		deleteConflicts["countries"] = countries
		deleteConflicts["home-slider"] = slides
		deleteConflicts["vendors"] = vendors
		deleteConflicts["news"] = news
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

func DeleteFile(c *gin.Context) {
	var file model.File
	if err := c.Bind(&file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Delete(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	path := os.Getenv("IMAGES_PATH") + file.Path
	if err := os.Remove(path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": file.ID, "status": "deleted"})
}

const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func generateFilePath(basePath string) (string, error) {
	firstDir := randStringBytes(2)
	secondDir := randStringBytes(2)

	path := firstDir + "/" + secondDir

	if err := os.MkdirAll(basePath+path, os.ModePerm); err != nil {
		return "", err
	}

	return path + "/", nil
}
