package fixtures

import (
	"fmt"

	"go-challenge/internal/models"

	"github.com/jinzhu/gorm"
)

func NewAnnonceFixture(cat *models.Cats) *models.Annonce {
	titles := []string{"Annonce de chat", "Chat à adopter", "Nouveau chat disponible", "Chat mignon à adopter"}
	descriptions := []string{
		"Ce chat est très amical et joueur.",
		"Adoptez ce chat affectueux et doux.",
		"Ce chat est en bonne santé et prêt à être adopté.",
		"Un chat parfait pour une famille aimante.",
	}

	title := randomChoice(titles)
	description := randomChoice(descriptions)

	return &models.Annonce{
		Title:       title,
		Description: &description,
		UserID:      cat.UserID,
		CatID:       fmt.Sprintf("%d", cat.ID),
	}
}

func CreateAnnonceFixtures(db *gorm.DB, cats []*models.Cats) error {
	for _, cat := range cats {
		annonce := NewAnnonceFixture(cat)
		if err := db.Create(annonce).Error; err != nil {
			return err
		}
	}
	return nil
}
