package server

import (
	// "fmt"
	"fmt"
	"net/http"

	"go-challenge/internal/handlers"
	// "go-challenge/internal/models"

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

		if role != "ADMIN" {
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

func FeatureFlagMiddleware(featureFlagHandler *handlers.FeatureFlagHandler, featureFlagName string) func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            isEnabled, err := featureFlagHandler.IsFeatureFlagEnabled(featureFlagName)
            if err != nil {
                fmt.Println(err)
                fmt.Println("1", featureFlagName)
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            if !isEnabled {
                fmt.Println("2", featureFlagName)
                http.Error(w, "Feature is not enabled", http.StatusForbidden)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
