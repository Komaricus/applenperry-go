package api

import (
	"crypto/sha256"
	"fmt"
	"github.com/applenperry-go/db"
	"github.com/applenperry-go/model"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
)

func CreateAdmin (c *gin.Context) {
	var admin model.Admin
	if err := c.Bind(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var check model.Admin
	//check if login exists
	if result := db.DB.Where("login = ?", admin.Login).First(&check); result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		return
	}

	//generate uuid id
	id, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	admin.ID = id.String()

	//generated hashed password
	h := sha256.New()
	h.Write([]byte(admin.Password + admin.ID))
	h.Write(h.Sum(nil))
	admin.Password = fmt.Sprintf("%x", h.Sum(nil))

	//store to db
	if err := db.DB.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, admin)
}