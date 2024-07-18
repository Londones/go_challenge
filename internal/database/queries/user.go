package queries

import (
	"go-challenge/internal/models"
	"go-challenge/internal/utils"
)

func (s *DatabaseService) FindUserByEmail(email string) (*models.User, error) {
	db := s.s.DB()
	var user models.User
	if err := db.Preload("Roles").Where("email = ?", email).First(&user).Error; err != nil {
		utils.Logger("error", "Find User By Email:", "Failed to find user by email", err.Error())
		return nil, err
	}
	return &user, nil
}

func (s *DatabaseService) FindUserByID(id string) (*models.User, error) {
	db := s.s.DB()
	var user models.User
	if err := db.Preload("Roles").Where("ID = ?", id).First(&user).Error; err != nil {
		utils.Logger("error", "Find User By ID:", "Failed to find user by ID", err.Error())
		return nil, err
	}

	return &user, nil
}

func (s *DatabaseService) FindUserByGoogleID(googleID string) (*models.User, error) {
	db := s.s.DB()
	var user models.User
	if err := db.Where("google_ID = ?", googleID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *DatabaseService) CreateUser(user *models.User, role *models.Roles) error {
	db := s.s.DB()
	if err := db.Create(user).Error; err != nil {
		return err
	}
	if err := db.Model(user).Association("Roles").Append(role).Error; err != nil {
		return err
	}
	return nil
}

func (s *DatabaseService) GetUserFavorites(UserID string) ([]models.Favorite, error) {
	db := s.s.DB()
	var favorites []models.Favorite
	if err := db.Where("UserID = ?", UserID).Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}

func (s *DatabaseService) UpdateUser(user *models.User) error {
	db := s.s.DB()
	if err := db.Save(user).Error; err != nil {
		return err
	}

	association := db.Model(user).Association("Roles").Replace(user.Roles)
	if association.Error != nil {
		return association.Error
	}
	return nil
}

func (s *DatabaseService) DeleteUser(id string) error {
	db := s.s.DB()
	var user models.User
	if err := db.Where("ID = ?", id).First(&user).Error; err != nil {
		return err
	}
	if err := db.Delete(&user).Error; err != nil {
		return err
	}
	return nil
}

func (s *DatabaseService) GetAllUsers() ([]models.User, error) {
	db := s.s.DB()
	var users []models.User
	if err := db.Preload("Roles").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *DatabaseService) GetRoleByName(name models.RoleName) (*models.Roles, error) {
	db := s.s.DB()
	var role models.Roles
	if err := db.Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
