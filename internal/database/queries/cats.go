package queries

import (
	"go-challenge/internal/models"
)

func (s *DatabaseService) CreateCat(cat *models.Cats) (id uint, err error) {
	db := s.s.DB()
	if err := db.Create(cat).Error; err != nil {
		return 0, err
	}
	return cat.ID, nil
}

func (s *DatabaseService) FindCatByID(id string) (*models.Cats, error) {
	db := s.s.DB()
	var cat models.Cats
	if err := db.Where("ID = ?", id).First(&cat).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}
