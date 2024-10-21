package services

import (
	"messaging-service/database"
	"messaging-service/models"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(username, password string) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{Username: username, Password: string(hashedPassword)}
	return database.DB.Create(&user).Error
}

func GetUser(username string) (models.User, error) {
	var user models.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	return user, err
}

func ListUsers() ([]string, error) {
	var users []string
	if err := database.DB.Model(&models.User{}).Pluck("username", &users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func ValidateUser(username, password string) bool {
	user, err := GetUser(username)
	if err != nil {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil
}
