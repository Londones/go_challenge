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

func CreateAssociation(association *models.Association) error {
	return db.Create(association).Error
}
