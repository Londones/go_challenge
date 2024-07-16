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

func (s *DatabaseService) CreateAssociation(association *models.Association) (uint, error) {
	db := s.s.DB()
	if err := db.Create(&association).Error; err != nil {
		return 0, err
	}
	return association.ID, nil
}

func (s *DatabaseService) GetAllAssociations() ([]models.Association, error) {
	db := s.s.DB()
	var associations []models.Association
	if err := db.Preload("Owner").Order("verified ASC").Find(&associations).Error; err != nil {
		return nil, err
	}

	return associations, nil
}

func (s *DatabaseService) UpdateAssociation(association *models.Association) error {
	db := s.s.DB()
	if err := db.Save(association).Error; err != nil {
		return err
	}
	return nil
}

func (s *DatabaseService) FindAssociationById(id int) (*models.Association, error) {
	db := s.s.DB()
	var association models.Association
	if err := db.Preload("Owner").First(&association, id).Error; err != nil {
		return nil, err
	}
	return &association, nil
}

func (s *DatabaseService) FindAssociationsByUserId(userId string) ([]models.Association, error) {
	db := s.s.DB()
	var associations []models.Association
	if err := db.Preload("Owner").Where("owner_id = ?", userId).Find(&associations).Error; err != nil {
		return nil, err
	}
	return associations, nil
}

func (s *DatabaseService) DeleteAssociation(id int) error {
	db := s.s.DB()
	if err := db.Delete(&models.Association{}, id).Error; err != nil {
		return err
	}
	return nil
}
