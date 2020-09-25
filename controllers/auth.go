package controllers

import (
	"net/http"
	"time"

	"gika.test/v1/helper"
	"gika.test/v1/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("halomamangGarox")

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func LoginController(c *gin.Context) {
	var input LoginInput
	var user models.Users
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.DB.Where("username=?", input.Username).First(&user).Error; err != nil {
		helper.RespondWithError(c, 401, "oops!!")
		return
	}

	match := CheckPasswordHash(input.Password, user.Password)

	if !match {
		helper.RespondWithError(c, 401, "credential anda salah")
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: input.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Model(&user).Update("token", tokenString)

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "username": input.Username})
}

func SingUpController(c *gin.Context) {
	var input LoginInput
	var user models.Users
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.DB.Where("username=?", input.Username).First(&user).Error; err != nil {

		expirationTime := time.Now().Add(5 * time.Minute)
		hashPwd, _ := HashPassword(input.Password)
		claims := &Claims{
			Username: input.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user := models.Users{Username: input.Username, Password: hashPwd, Token: tokenString}
		models.DB.Create(&user)
		c.JSON(http.StatusOK, gin.H{"token": tokenString, "username": input.Username})
		return
	}
	helper.RespondWithError(c, 404, "username sudah dipakai")

}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
