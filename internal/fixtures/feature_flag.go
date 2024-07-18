package fixtures

import (
	"go-challenge/internal/models"

	"github.com/jinzhu/gorm"
)

func NewFeatureFlagFixture(index int) *models.FeatureFlag {
    names := []string{"OAuth", "Association"}

    return &models.FeatureFlag{
        Name:      names[index],
        IsEnabled: true,
    }
}

func CreateFeatureFlagFixture(db *gorm.DB) error {
    for i := 0; i < 2; i++ {
        featureFlag := NewFeatureFlagFixture(i)
        if err := db.Create(featureFlag).Error; err != nil {
            return err
        }
    }
    return nil
}