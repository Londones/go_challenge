package queries

import (
	"go-challenge/internal/models"
)

func (s *DatabaseService) CreateFavorite(favorite *models.Favorite) (uint, error) {
	db := s.s.DB()
	if err := db.Create(favorite).Error; err != nil {
		return 0, err
	}

	return favorite.ID, nil
}

func (s *DatabaseService) UpdateFavorite(favorite *models.Favorite) error {
	db := s.s.DB()
	return db.Save(favorite).Error
}

func (s *DatabaseService) DeleteFavorite(favorite *models.Favorite) error {
	db := s.s.DB()
	return db.Delete(favorite).Error
}

func (s *DatabaseService) FindFavoriteByID(id string) (*models.Favorite, error) {
	db := s.s.DB()
	var favorite models.Favorite
	if err := db.Where("ID = ?", id).First(&favorite).Error; err != nil {
		return nil, err
	}
	return &favorite, nil
}

func (s *DatabaseService) FindFavoritesByUserID(userID string) ([]models.Favorite, error) {
	db := s.s.DB()
	var favorites []models.Favorite
	if err := db.Where("user_id = ?", userID).Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}

func (s *DatabaseService) FindFavoritesByAnnonceID(annonceID string) ([]models.Favorite, error) {
	db := s.s.DB()
	var favorites []models.Favorite
	if err := db.Where("annonce_id = ?", annonceID).Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}

func (s *DatabaseService) FindFavoriteByUserAndAnnonceID(userID, annonceID string) (*models.Favorite, error) {
	db := s.s.DB()
	var favorite models.Favorite
	if err := db.Where("user_id = ? AND annonce_id = ?", userID, annonceID).First(&favorite).Error; err != nil {
		return nil, err
	}
	return &favorite, nil
}
