package main

import (
	"net/http"

	"gika.test/v1/controllers"
	"gika.test/v1/middleware"
	"gika.test/v1/models"
	"github.com/axiaoxin-com/ratelimiter"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})
	models.ConnectDatabase()
	router.Use(ratelimiter.GinMemRatelimiter(1000*1000, 2))
	router.POST("/masuk", controllers.LoginController)
	router.POST("/daftar", controllers.SingUpController)
	auth := router.Group("/auth")
	auth.Use(middleware.TokenAuthMiddleware())
	{
		auth.GET("/books", controllers.AllBooks)
		auth.GET("/pengguna", controllers.GetPengguna)
		auth.POST("/books", controllers.CreateBook)
	}
	router.Run(":8081")
}
