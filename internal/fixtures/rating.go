package fixtures

import (
	"go-challenge/internal/models"
	"math/rand"

	"github.com/jinzhu/gorm"
)

func NewRatingFixture(userID string, authorID string) *models.Rating {
	comments := []string{"Great!", "Not bad", "Could be better", "Excellent!", "Just okay"}

	return &models.Rating{
		Mark:     int8(rand.Intn(5) + 1),
		Comment:  randomChoice(comments),
		UserID:   userID,
		AuthorID: authorID,
	}
}

func CreateRatingFixtures(db *gorm.DB, userID string, authorID string, count int) error {
	for i := 0; i < count; i++ {
		rating := NewRatingFixture(userID, authorID)
		if err := db.Create(rating).Error; err != nil {
			return err
		}
	}
	return nil
}
