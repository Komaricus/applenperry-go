package api

import (
	"fmt"
	"github.com/applenperry-go/config"
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

func UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	files := form.File["files"]

	dbFiles := make([]model.File, 0)
	var basePath string
	configuration := config.GetConfig()
	if configuration.PRODUCTION {
		//todo
		basePath = "/"
	} else {
		basePath = "/Users/aleksandr/Documents/applenperry/applenperry-vue/src/assets/img/"
	}

	for _, file := range files {
		id, err := uuid.NewV4()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		originalName := filepath.Base(file.Filename)
		ext := filepath.Ext(file.Filename)
		size := file.Size
		filename := id.String() + ext

		filePath, err := generateFilePath(basePath)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		if err := c.SaveUploadedFile(file, basePath+filePath+filename); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		dbFile := model.File{
			ID:           id.String(),
			FileName:     filename,
			Path:         filePath + filename,
			OriginalName: originalName,
			Size:         size,
		}

		if err := db.DB.Create(&dbFile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		dbFiles = append(dbFiles, dbFile)
	}
	c.JSON(http.StatusOK, dbFiles)
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
