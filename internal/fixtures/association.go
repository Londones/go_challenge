package fixtures

import (
	"go-challenge/internal/models"

	"github.com/jinzhu/gorm"
)

func NewAssociationFixture(ownerID string) *models.Association {
	var associationIsVerified = true

	return &models.Association{
		Name:       "Assoc de TEST",
		AddressRue: "10 rue du TEST",
		Cp:         "12345",
		Ville:      "TEST-CITY",
		Phone:      "0101010101",
		Email:      "emailDeTest@gmail.com",
		OwnerID:    ownerID,
		Verified:   &associationIsVerified,
	}
}

func CreateAssociationFixtures(db *gorm.DB, userID string) (*models.Association, error) {
	association := NewAssociationFixture(userID)
	if err := db.Create(&association).Error; err != nil {
		return nil, err
	}
	return association, nil
}

func AssociationAddMembers(db *gorm.DB, asso *models.Association, users []models.User) (*models.Association, error) {
	asso.Members = users
	if err := db.Save(asso).Error; err != nil {
		return nil, err
	}
	return asso, nil
}
