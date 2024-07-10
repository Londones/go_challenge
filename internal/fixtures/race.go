package fixtures

import (
	"go-challenge/internal/models"

	"github.com/jinzhu/gorm"
)

func NewRaceFixture(index int) *models.Races {
	names := []string{"American shorthair", "Munchkin", "Main coon", "Toyger", "Ragdoll"}

	return &models.Races{
		RaceName: names[index],
	}
}

func CreateRaceFixture(db *gorm.DB) error {
	for i := 0; i < 5; i++ {
		race := NewRaceFixture(i)
		if err := db.Create(race).Error; err != nil {
			return err
		}
	}
	return nil
}
