package queries

import (
	"go-challenge/internal/models"
	"go-challenge/internal/utils"
)

func (s *DatabaseService) GetAllFeatureFlags() ([]models.FeatureFlag, error) {
	db := s.s.DB()
	var featureFlags []models.FeatureFlag
	if err := db.Find(&featureFlags).Error; err != nil {
		utils.Logger("error", "Get Feature Flags:", "Failed to get feature flags", err.Error())
		return nil, err
	}
	return featureFlags, nil
}

func (s *DatabaseService) FindFeatureFlagByID(id int) (models.FeatureFlag, error) {
	db := s.s.DB()
	var featureFlag models.FeatureFlag
	if err := db.Where("id = ?", id).First(&featureFlag).Error; err != nil {
		utils.Logger("error", "Find Feature Flag By ID:", "Failed to find feature flag by ID", err.Error())
		return featureFlag, err
	}
	return featureFlag, nil
}

func (s *DatabaseService) UpdateFeatureFlag(featureFlag models.FeatureFlag) error {
	db := s.s.DB()
	if err := db.Save(featureFlag).Error; err != nil {
		utils.Logger("error", "Update Feature Flag:", "Failed to update feature flag", err.Error())
		return err
	}
	return nil
}

func (s *DatabaseService) FindFeatureFlagByName(name string) (models.FeatureFlag, error) {
	db := s.s.DB()
	var featureFlag models.FeatureFlag
	if err := db.Where("name = ?", name).First(&featureFlag).Error; err != nil {
		utils.Logger("error", "Find Feature Flag By Name:", "Failed to find feature flag by name", err.Error())
		return featureFlag, err
	}
	return featureFlag, nil
}