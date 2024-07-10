package fixtures

import (
	"fmt"
	"log"
	"time"

	"go-challenge/internal/auth"
	"go-challenge/internal/models"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

func NewUserFixture(i int) *models.User {
	name := fmt.Sprintf("user%d", i)
	email := fmt.Sprintf("user%d@gmail.com", i)
	password, err := auth.HashPassword("password")
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}
	addressRues := []string{"123 Main St", "456 Elm St", "789 Oak St", "101 Maple St", "202 Pine St"}
	cps := []string{"75001", "75002", "75003", "75004", "75005"}
	villes := []string{"Paris", "Lyon", "Marseille", "Toulouse", "Nice"}

	return &models.User{
		ID:            uuid.New().String(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Name:          name,
		Email:         email,
		Password:      password,
		AddressRue:    randomChoice(addressRues),
		Cp:            randomChoice(cps),
		Ville:         randomChoice(villes),
		Associations:  []models.Association{},
		Roles:         []models.Roles{},
		GoogleID:      "",
		ProfilePicURL: "default",
	}
}

func CreateUserFixtures(db *gorm.DB, count int, userRole *models.Roles) ([]*models.User, error) {
	var users []*models.User

	for i := 1; i <= count; i++ {
		user := NewUserFixture(i)
		if err := db.Create(user).Error; err != nil {
			return nil, fmt.Errorf("failed to create user %d: %v", i, err)
		}
		// Assign the user role to the newly created user
		if err := db.Model(user).Association("Roles").Append(userRole).Error; err != nil {
			return nil, fmt.Errorf("failed to assign role to user %d: %v", i, err)
		}
		users = append(users, user)
	}
	return users, nil
}
