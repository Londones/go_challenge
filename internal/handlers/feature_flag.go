package handlers

import (
	"encoding/json"
	"fmt"
	"go-challenge/internal/database/queries"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type FeatureFlagHandler struct {
	featureFlagQueries *queries.DatabaseService
}

func NewFeatureFlagHandler(featureFlagQueries *queries.DatabaseService) *FeatureFlagHandler {
	return &FeatureFlagHandler{featureFlagQueries: featureFlagQueries}
}

// GetAllFeatureFlagsHandler godoc
// @Summary Get all feature flags
// @Description Get all feature flags
// @Tags featureFlags
// @Accept  json
// @Produce  json
// @Success 200 {array} models.FeatureFlag
// @Failure 500 {string} string "Internal server error"
// @Router /featureflags [get]
func (h *FeatureFlagHandler) GetAllFeatureFlagsHandler(w http.ResponseWriter, r *http.Request) {
	featureFlags, err := h.featureFlagQueries.GetAllFeatureFlags()
	if err != nil {
		http.Error(w, "error fetching feature flags", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(featureFlags)
	if err != nil {
		http.Error(w, "error encoding feature flags to JSON", http.StatusInternalServerError)
		return
	}
}

// UpdateFeatureFlagStatusHandler godoc
// @Summary Update feature flag status
// @Description Update the status of a feature flag
// @Tags featureFlags
// @Accept  json
// @Produce  json
// @Param id path int true "Feature Flag ID"
// @Param isEnabled body bool true "Is Enabled"
// @Success 200 {object} models.FeatureFlag
// @Failure 400 {object} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /featureflags/{id} [put]
func (h *FeatureFlagHandler) UpdateFeatureFlagStatusHandler(w http.ResponseWriter, r *http.Request) {
	featureFlagIDStr := chi.URLParam(r, "id")
	if featureFlagIDStr == "" {
		http.Error(w, "missing feature flag ID", http.StatusBadRequest)
		return
	}

	featureFlagID, err := strconv.Atoi(featureFlagIDStr)
	if err != nil {
		http.Error(w, "invalid feature flag ID", http.StatusBadRequest)
		return
	}

	var body struct {
		IsEnabled bool `json:"isEnabled"`
	}

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	featureFlag, err := h.featureFlagQueries.FindFeatureFlagByID(featureFlagID)
	if err != nil {
		http.Error(w, "error fetching feature flag", http.StatusInternalServerError)
		return
	}

	featureFlag.IsEnabled = body.IsEnabled

	err = h.featureFlagQueries.UpdateFeatureFlag(featureFlag)
	if err != nil {
		http.Error(w, "error updating feature flag", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(featureFlag)
}

func (h *FeatureFlagHandler) IsFeatureFlagEnabled(featureFlagName string) (bool, error) {
	featureFlag, err := h.featureFlagQueries.FindFeatureFlagByName(featureFlagName)
	if err != nil {
		fmt.Println("error fetching feature flag: %w", err)
		fmt.Println("3", featureFlagName)
		return false, fmt.Errorf("error fetching feature flag: %w", err)
	}

	log.Printf("Feature flag '%s' is enabled: %v", featureFlagName, featureFlag.IsEnabled)
	return featureFlag.IsEnabled, nil
}
