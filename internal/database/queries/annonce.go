package queries

import (
	"errors"
	"fmt"
	"go-challenge/internal/models"

	"gorm.io/gorm"
)

func (s *DatabaseService) CreateAnnonce(annonce *models.Annonce) (id uint, err error) {
	db := s.s.DB()
	if err := db.Create(annonce).Error; err != nil {
		return 0, err
	}
	return annonce.ID, nil
}

func (s *DatabaseService) FindAnnonceByID(id string) (*models.Annonce, error) {
	db := s.s.DB()
	var annonce models.Annonce
	if err := db.Where("ID = ?", id).First(&annonce).Error; err != nil {
		return nil, err
	}
	return &annonce, nil
}

func (s *DatabaseService) UpdateAnnonceDescription(id string, description string) (*models.Annonce, error) {
	db := s.s.DB()

	var annonce models.Annonce
	if err := db.Where("id = ?", id).First(&annonce).Error; err != nil {
		return nil, err
	}

	annonce.Description = &description

	if err := db.Save(&annonce).Error; err != nil {
		return nil, err
	}

	return &annonce, nil
}

func (s *DatabaseService) DeleteAnnonce(id string) error {
	db := s.s.DB()

	// Vérifier si l'annonce existe avant de tenter de la supprimer
	var annonce models.Annonce
	if err := db.Where("id = ?", id).First(&annonce).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("annonce with ID %s not found", id)
		}
		return err
	}

	if err := db.Delete(&annonce).Error; err != nil {
		return err
	}

	return nil
}

func (s *DatabaseService) GetAllAnnonces() ([]models.Annonce, error) {
	db := s.s.DB()
	var annonces []models.Annonce
	if err := db.Find(&annonces).Error; err != nil {
		return nil, err
	}
	return annonces, nil
}

func (s *DatabaseService) GetUserAnnonces(userID string) ([]models.Annonce, error) {
	db := s.s.DB()
	var annonces []models.Annonce
	if err := db.Where("user_id = ?", userID).Find(&annonces).Error; err != nil {
		return nil, err
	}
	return annonces, nil
}

// get l'annonce d'un chat par son id
func (s *DatabaseService) FindAnnonceByCatID(catID string) (*models.Annonce, error) {
	db := s.s.DB()
	var annonce models.Annonce
	if err := db.Where("cat_id = ?", catID).First(&annonce).Error; err != nil {
		return nil, err
	}
	return &annonce, nil
}

func (s *DatabaseService) GetUserIDByAnnonceID(annonceID string) (id string, err error) {
	db := s.s.DB()
	var user models.User

	if err := db.Table("users").Joins("JOIN annonces ON users.id = annonces.user_id::uuid").Where("annonces.id = ?", annonceID).First(&user).Error; err != nil {
		return "", err
	}
	return user.ID, nil
}

func (s *DatabaseService) GetAddressFromAnnonceID(userID string) (string, error) {
	db := s.s.DB()
	var user models.User
	if err := db.Where("ID = ?", userID).First(&user).Error; err != nil {
		return "", err
	}
	var address = user.AddressRue + user.Cp + user.Ville
	fmt.Println(address)
	return address, nil
}
