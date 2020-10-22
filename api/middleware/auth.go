package middleware

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/applenperry-go/db"
	"github.com/applenperry-go/model"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type login struct {
	Login    string `form:"login" json:"login" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var identityKey = "id"

type User struct {
	ID   string
	Role string
}

func verifyCredential(admin login) (*User, error) {
	var fromDB model.Admin
	//get user from db by login
	if err := db.DB.Where("login = ?", admin.Login).First(&fromDB).Error; err != nil {
		return nil, err
	}

	//generated hashed password
	h := sha256.New()
	h.Write([]byte(admin.Password + fromDB.ID))
	h.Write(h.Sum(nil))
	admin.Password = fmt.Sprintf("%x", h.Sum(nil))

	//compare passwords
	if admin.Password != fromDB.Password {
		return nil, errors.New("wrong password")
	}

	return &User{
		ID:   fromDB.ID,
		Role: "admin",
	}, nil
}

func GetAuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "applenperry",
		Key:         []byte(os.Getenv("ACCESS_SECRET")),
		Timeout:     time.Hour,
		MaxRefresh:  8 * time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
					"role":      v.Role,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				ID:   claims[identityKey].(string),
				Role: claims["role"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", errors.New("missing login or password")
			}

			user, err := verifyCredential(loginVals)
			if err == nil {
				return user, nil
			}

			return nil, errors.New("incorrect login or password")
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*User); ok && v.Role == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		return nil, errors.New("JWT Error:" + err.Error())
	}

	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		return nil, errors.New("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	return authMiddleware, nil
}
