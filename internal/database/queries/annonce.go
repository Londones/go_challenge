package queries

import (
	"go-challenge/internal/models"
)

func (s *DatabaseService) CreateAnnonce(annonce *models.Annonce) (id uint, err error) {
	db := s.s.DB()
	if err := db.Create(annonce).Error; err != nil {
		return 0, err
	}
	return annonce.ID, nil
}
