package server

import (
	"fmt"
	"net/http"

	"go-challenge/internal/handlers"
	"go-challenge/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

type Middleware struct {
	FeatureFlagHandler *handlers.FeatureFlagHandler
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		role := claims["role"]

		if role != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func UserOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil || claims["id"] == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID := claims["id"].(string)
		routeUserID := chi.URLParam(r, "id")

		// Only check the route user ID if it exists
		if routeUserID != "" && userID != routeUserID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) FeatureFlagMiddleware(featureName string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		featureFlags, err := m.FeatureFlagHandler.GetFeatureFlags()
		if err != nil {
			http.Error(w, "Error fetching feature flags", http.StatusInternalServerError)
			return
		}

		enabled, err := IsFeatureEnabled(featureFlags, featureName)
		if err != nil {
			http.Error(w, "Error checking feature flag", http.StatusInternalServerError)
			return
		}
		if !enabled {
			http.Error(w, "Feature not enabled", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func IsFeatureEnabled(featureFlags []models.FeatureFlag, featureName string) (bool, error) {
	for _, featureFlag := range featureFlags {
		if featureFlag.Name == featureName {
			return featureFlag.IsEnabled, nil
		}
	}

	return false, fmt.Errorf("feature flag not found: %s", featureName)
}
