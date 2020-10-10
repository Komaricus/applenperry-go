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
	"time"
)

func GetFiles(c *gin.Context) {
	var files []model.File

	if err := db.DB.Where("is_deleted = false").Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, files)
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

	if len(countries) > 0 || len(slides) > 0 {
		c.JSON(http.StatusOK, gin.H{"id": id, "status": "not_deletable", "countries": countries, "homeSlides": slides})
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

	path := os.Getenv("IMAGES_PATH") + file.Path
	if err := os.Remove(path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Delete(&file).Error; err != nil {
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
