package queries

import (
	"go-challenge/internal/database"
	"go-challenge/internal/models"
)

type DatabaseService struct {
	s database.Service
}

func NewQueriesService(s *database.Service) *DatabaseService {
	return &DatabaseService{
		s: *s,
	}
}

func (s *DatabaseService) CreateAssociation(association *models.Association) error {
	db := s.s.DB()
	if err := db.Create(association).Error; err != nil {
		return err
	}
	return nil
}
