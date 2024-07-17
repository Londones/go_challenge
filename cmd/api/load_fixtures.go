package main

import (
	"fmt"
	"log"
	"os"

	"go-challenge/internal/fixtures"
	"go-challenge/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func ExecuteFixtures() {
	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	// Load the fixtures.
	LoadFixtures(db)
}

func LoadFixtures(db *gorm.DB) {
	// Get the USER role
	var userRole models.Roles
	if err := db.Where("name = ?", models.UserRole).First(&userRole).Error; err != nil {
		fmt.Printf("failed to find user role: %v", err)
	}

	// Create 5 races
	err := fixtures.CreateRaceFixture(db)
	if err != nil {
		fmt.Printf("failed to create race fixture: %v", err)
	}

	// Create 5 users
	users, err := fixtures.CreateUserFixtures(db, 5, &userRole)
	if err != nil {
		fmt.Printf("failed to create user fixtures: %v", err)
	}

	// For each user, create 5 cats and 5 corresponding annonces
	for _, user := range users {
		cats, err := fixtures.CreateCatFixturesForUser(db, 5, user.ID)
		if err != nil {
			fmt.Printf("failed to create cat fixtures for user %s: %v", user.ID, err)
		}

		if err := fixtures.CreateAnnonceFixtures(db, cats); err != nil {
			fmt.Printf("failed to create annonce fixtures for user %s: %v", user.ID, err)
		}
	}

	// Création des fixtures pour les évaluations
	staticUserID := "38f5ca5d-0c87-425f-97fe-c84c3ee0997c"
	staticAuthorID := "5a7a8b69-6f8d-4818-ac15-b6a83b4fe518"
	err = fixtures.CreateRatingFixtures(db, staticUserID, staticAuthorID, 5)
	if err != nil {
		fmt.Printf("failed to create rating fixtures: %v", err)
	}
}
