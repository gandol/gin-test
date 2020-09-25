package helper

import (
	"gika.test/v1/models"
)

func GetDataUsers(token string) models.Users {
	var users models.Users
	if err := models.DB.Where("token=?", token).First(&users).Error; err != nil {
		return users
	}
	return users
}
