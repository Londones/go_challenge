package queries

import (
	"go-challenge/internal/models"
	"go-challenge/internal/utils"
)

func (s *DatabaseService) GetFeatureFlags() ([]models.FeatureFlag, error) {
	db := s.s.DB()
	var featureFlags []models.FeatureFlag
	if err := db.Find(&featureFlags).Error; err != nil {
		utils.Logger("error", "Get Feature Flags:", "Failed to get feature flags", err.Error())
		return nil, err
	}
	return featureFlags, nil
}

func (s *DatabaseService) UpdateFeatureFlag(featureFlag models.FeatureFlag) error {
	db := s.s.DB()
	if err := db.Save(featureFlag).Error; err != nil {
		utils.Logger("error", "Update Feature Flag:", "Failed to update feature flag", err.Error())
		return err
	}
	return nil
}
