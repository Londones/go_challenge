package queries

import (
	"go-challenge/internal/database"
	"go-challenge/internal/models"
)

type UserQueries interface {
	FindUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
}

var (
	s  = database.Service{}
	db = s.DB()
)

func FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *models.User) error {
	return db.Create(user).Error
}
