package handlers

import (
	"go-challenge/internal/database/queries"
	"go-challenge/internal/models"
)

type FeatureFlagHandler struct {
	featureFlagQueries *queries.DatabaseService
}

func NewFeatureFlagHandler(featureFlagQueries *queries.DatabaseService) *FeatureFlagHandler {
	return &FeatureFlagHandler{featureFlagQueries: featureFlagQueries}
}

func (h *FeatureFlagHandler) GetFeatureFlags() ([]models.FeatureFlag, error) {
	featureFlags, err := h.featureFlagQueries.GetFeatureFlags()
	if err != nil {
		return nil, err
	}

	featureFlagCache := make(map[string]bool)
	for _, featureFlag := range featureFlags {
		featureFlagCache[featureFlag.Name] = featureFlag.IsEnabled
	}

	return featureFlags, nil
}

func (h *FeatureFlagHandler) UpdateFeatureFlag(featureFlag models.FeatureFlag) error {
	return h.featureFlagQueries.UpdateFeatureFlag(featureFlag)
}
