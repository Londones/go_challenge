package queries

import (
	"go-challenge/internal/models"
)

func (s *DatabaseService) FindUserByEmail(email string) (*models.User, error) {
	db := s.s.DB()
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *DatabaseService) FindUserByID(id string) (*models.User, error) {
	db := s.s.DB()
	var user models.User
	if err := db.Where("ID = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *DatabaseService) FindUserByGoogleID(googleID string) (*models.User, error) {
	db := s.s.DB()
	var user models.User
	if err := db.Where("googleID = ?", googleID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *DatabaseService) CreateUser(user *models.User) error {
	db := s.s.DB()
	return db.Create(user).Error
}

func (s *DatabaseService) GetUserFavorites(userID string) ([]models.Favorite, error) {
	db := s.s.DB()
	var favorites []models.Favorite
	if err := db.Where("userID = ?", userID).Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}

func (s *DatabaseService) UpdateUser(user *models.User) error {
	db := s.s.DB()
	return db.Save(user).Error
}

func (s *DatabaseService) DeleteUser(id string) error {
	db := s.s.DB()
	return db.Where("id = ?", id).Delete(&models.User{}).Error
}

func (s *DatabaseService) GetAllUsers() ([]models.User, error) {
	db := s.s.DB()
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
