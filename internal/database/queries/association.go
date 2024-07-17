package queries

import (
	"go-challenge/internal/database"
	"go-challenge/internal/models"

	"github.com/jinzhu/gorm"
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

func (s *DatabaseService) GetAllAssociations() ([]models.Association, error) {
	db := s.s.DB()
	var associations []models.Association
	if err := db.Order("verified ASC").Find(&associations).Error; err != nil {
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
	if err := db.First(&association, id).Error; err != nil {
		return nil, err
	}
	return &association, nil
}

func (s *DatabaseService) FindAssociationsByUserId(userId string) ([]models.Association, error) {
	db := s.s.DB()
	var associations []models.Association
	if err := db.Where("owner_id = ? OR ? = ANY(members)", userId, userId).Find(&associations).Error; err != nil {
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

func (s *DatabaseService) UpdateAssociationMembers(associationID uint, memberIDs []string) error {
	db := s.s.DB()
	tx := db.Begin()

	if err := tx.Model(&models.Association{Model: gorm.Model{ID: associationID}}).Association("Members").Clear().Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, memberID := range memberIDs {
		user := models.User{}
		if err := tx.Where("id = ?", memberID).First(&user).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Model(&models.Association{Model: gorm.Model{ID: associationID}}).Association("Members").Append(&user).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
