package middleware

import (
	"log"

	"gika.test/v1/helper"
	"gika.test/v1/models"

	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	requiredToken := "aaa"

	// We want to make sure the token is set, bail if not
	if requiredToken == "" {
		log.Fatal("Please set API_TOKEN environment variable")
	}
	return func(c *gin.Context) {
		var users models.Users
		token := c.Request.Header.Get("token")
		if token == "" {
			helper.RespondWithError(c, 401, "not allowed")
			return
		}
		if err := models.DB.Where("token=?", token).First(&users).Error; err != nil {
			helper.RespondWithError(c, 401, "Invalid Credential")
			return
		}
		c.Set("userId", users.ID)
		c.Set("username", users.Username)
		c.Next()
	}
}
