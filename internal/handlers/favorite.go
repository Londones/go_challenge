package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"go-challenge/internal/database/queries"
	"go-challenge/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jinzhu/gorm"
)

type FavoriteHandler struct {
	favoriteQueries *queries.DatabaseService
	userQueries     *queries.DatabaseService
}

func NewFavoriteHandler(favoriteQueries, userQueries *queries.DatabaseService) *FavoriteHandler {
	return &FavoriteHandler{favoriteQueries: favoriteQueries, userQueries: userQueries}
}

// FavoriteCreationHandler godoc
// @Summary Create favorites
// @Description Create a new favorite with the provided details
// @Tags favorites
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param annonceID formData string true "ID of the annonce"
// @Success 201 {object} models.Favorite "favorite created successfully"
// @Failure 400 {string} string "annonceID is required"
// @Failure 500 {string} string "error creating favorite"
// @Router /favorites [post]
func (h *FavoriteHandler) FavoriteCreationHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	var annonceID string

	if strings.Contains(contentType, "application/json") {
		var data struct {
			AnnonceID string `json:"annonceID"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		annonceID = data.AnnonceID
	} else {
		r.ParseForm()
		annonceID = r.FormValue("annonceID")
	}

	print(annonceID)

	if annonceID == "" {
		http.Error(w, "annonceID is required", http.StatusBadRequest)
		return
	}

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "error getting claims", http.StatusInternalServerError)
		return
	}
	userID := claims["id"].(string)

	user, err := h.userQueries.FindUserByID(userID)
	if err != nil {
		http.Error(w, "error finding user", http.StatusInternalServerError)
		return
	}

	_, err = h.favoriteQueries.FindAnnonceByID(annonceID)
	if err != nil {
		http.Error(w, "error finding annonce", http.StatusInternalServerError)
		return
	}

	favorite := &models.Favorite{
		UserID:    user.ID,
		AnnonceID: annonceID,
	}

	err = h.favoriteQueries.CreateFavorite(favorite)
	if err != nil {
		http.Error(w, "error creating favorite", http.StatusInternalServerError)
		return
	}

	authorID, err := h.userQueries.GetUserIDByAnnonceID(annonceID)
	if err != nil {
		http.Error(w, "error getting author ID", http.StatusInternalServerError)
		return
	}

	room := &models.Room{
		User1ID:   authorID,
		User2ID:   userID,
		AnnonceID: annonceID,
	}

	_, err = h.favoriteQueries.CreateRoom(room)
	if err != nil {
		http.Error(w, "error creating room", http.StatusInternalServerError)
		return
	}

	response := struct {
		Success  string           `json:"success"`
		Favorite *models.Favorite `json:"favorite"`
	}{
		Success:  "true",
		Favorite: favorite,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetFavoritesByUserHandler godoc
// @Summary Get user favorites
// @Description Get all favorites of the user
// @Tags favorites
// @Produce json
// @Param userID path string true "ID of the user"
// @Success 200 {array} models.Favorite "List of user favorites"
// @Failure 400 {string} string "user ID is required"
// @Failure 404 {string} string "favorites not found"
// @Failure 500 {string} string "error retrieving favorites"
// @Router /favorites/users/{userID} [get]
func (h *FavoriteHandler) GetFavoritesByUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	if userID == "" {
		http.Error(w, "user ID is required", http.StatusBadRequest)
		return
	}

	favorites, err := h.favoriteQueries.FindFavoritesByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, fmt.Sprintf("favorites for user with ID %s not found", userID), http.StatusNotFound)
			return
		}
		http.Error(w, "error retrieving favorites", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(favorites)
	if err != nil {
		http.Error(w, "error encoding favorites to JSON", http.StatusInternalServerError)
		return
	}
}

// DeleteFavoriteByIDHandler godoc
// @Summary Delete favorite by favorite ID
// @Description Delete a favorite by its ID
// @Tags favorites
// @Param favoriteID path string true "Favorite ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Favorite ID is required"
// @Failure 404 {string} string "Favorite not found"
// @Failure 500 {string} string "Error deleting favorite"
// @Router /favorites/{favoriteID} [delete]
func (h *FavoriteHandler) DeleteFavoriteByIDHandler(w http.ResponseWriter, r *http.Request) {
	favoriteID := chi.URLParam(r, "favoriteID")
	if favoriteID == "" {
		http.Error(w, "favorite ID is required", http.StatusBadRequest)
		return
	}

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "error getting claims", http.StatusInternalServerError)
		return
	}
	userID := claims["id"].(string)

	favorite, err := h.favoriteQueries.FindFavoriteByID(favoriteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, fmt.Sprintf("favorite with ID %s not found", favoriteID), http.StatusNotFound)
			return
		}
		http.Error(w, "error finding favorite", http.StatusInternalServerError)
		return
	}

	// Vérifiez si le favori appartient à l'utilisateur
	if favorite.UserID != userID {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.favoriteQueries.DeleteFavorite(favorite)
	if err != nil {
		http.Error(w, "error deleting favorite", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
