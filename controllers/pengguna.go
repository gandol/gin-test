package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPengguna(c *gin.Context) {
	user, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{"data": user})
}
