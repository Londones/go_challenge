package queries

import (
	"go-challenge/internal/models"
)

type UserQueries interface {
	FindUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
	GetUserFavorites(userID string) ([]models.Favorite, error)
	FindUserByID(id string) (*models.User, error)
}

func FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func FindUserByID(id string) (*models.User, error) {
	var user models.User
	if err := db.Where("ID = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func FindUserByGoogleID(googleID string) (*models.User, error) {
	var user models.User
	if err := db.Where("googleID = ?", googleID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *models.User) error {
	return db.Create(user).Error
}

func GetUserFavorites(userID string) ([]models.Favorite, error) {
	var favorites []models.Favorite
	if err := db.Where("userID = ?", userID).Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}
