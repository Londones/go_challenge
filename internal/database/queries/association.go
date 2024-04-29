package queries

import (
	"go-challenge/internal/database"
	"go-challenge/internal/models"
)

var (
	s  = database.Service{}
	db = s.DB()
)

type AssociationQueries interface {
	CreateAssociation(association *models.Association) error
}

func CreateAssociation(association *models.Association) (id uint, err error) {
	if err := db.Create(association).Error; err != nil {
		return 0, err
	}
	return association.ID, nil
}

func AddUserToAssociation(associationID uint, userID string) error {
	var association models.Association
	if err := db.Where("id = ?", associationID).First(&association).Error; err != nil {
		return err
	}

	association.MemberIDs = append(association.MemberIDs, userID)

	return db.Save(&association).Error
}
