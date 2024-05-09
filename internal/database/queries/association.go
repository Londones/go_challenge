package queries

import (
	"go-challenge/internal/database"
	"go-challenge/internal/models"
)

type DatabaseService struct {
	s database.Service
}

func NewQueriesService() *DatabaseService {
	return &DatabaseService{}
}

func (s *DatabaseService) CreateAssociation(association *models.Association) (id uint, err error) {
	db := s.s.DB()
	if err := db.Create(association).Error; err != nil {
		return 0, err
	}
	return association.ID, nil
}

func (s *DatabaseService) AddUserToAssociation(associationID uint, userID string) error {
	db := s.s.DB()
	var association models.Association
	if err := db.Where("id = ?", associationID).First(&association).Error; err != nil {
		return err
	}

	association.MemberIDs = append(association.MemberIDs, userID)

	return db.Save(&association).Error
}
