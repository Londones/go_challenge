package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"go-challenge/internal/database/queries"
	"go-challenge/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

type RatingHandler struct {
	ratingQueries *queries.DatabaseService
	userQueries   *queries.DatabaseService
}

func NewRatingHandler(ratingQueries, userQueries *queries.DatabaseService) *RatingHandler {
	return &RatingHandler{
		ratingQueries: ratingQueries,
		userQueries:   userQueries,
	}
}

// CreateRatingHandler handles the creation of a new rating.
// @Summary Create rating
// @Description Create a new rating with the provided details
// @Tags ratings
// @Accept x-www-form-urlencoded
// @Produce json
// @Param mark formData int true "Rating mark"
// @Param comment formData string false "Comment about the rating"
// @Param annonceID formData string true "Annonce ID"
// @Success 201 {object} models.Rating "Rating created successfully"
// @Failure 400 {string} string "Missing or invalid fields in the request"
// @Failure 500 {string} string "Internal server error"
// @Router /ratings [post]
func (h *RatingHandler) CreateRatingHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	markStr := r.FormValue("mark")
	comment := r.FormValue("comment")
	annonceID := r.FormValue("annonceID")

	if markStr == "" || annonceID == "" {
		http.Error(w, "Mark and annonceID are required", http.StatusBadRequest)
		return
	}

	mark, err := strconv.ParseInt(markStr, 10, 8)
	if err != nil {
		http.Error(w, "Invalid mark value", http.StatusBadRequest)
		return
	}

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "error getting claims", http.StatusInternalServerError)
		return
	}
	userID := claims["id"].(uint)

	rating := &models.Rating{
		Mark:      int8(mark),
		Comment:   comment,
		UserID:    userID,
		AnnonceID: annonceID,
	}

	ratingID, err := h.ratingQueries.CreateRating(rating)
	if err != nil {
		http.Error(w, "error creating rating", http.StatusInternalServerError)
		return
	}

	createdRating, err := h.ratingQueries.FindRatingByID(fmt.Sprintf("%d", ratingID))
	if err != nil {
		http.Error(w, "error retrieving created rating", http.StatusInternalServerError)
		return
	}

	response := struct {
		Success string         `json:"success"`
		Rating  *models.Rating `json:"rating"`
	}{
		Success: "true",
		Rating:  createdRating,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateRatingHandler handles updates to an existing rating.
// @Summary Update rating
// @Description Update the details of an existing rating
// @Tags ratings
// @Accept x-www-form-urlencoded
// @Produce json
// @Param id path string true "ID of the rating to update"
// @Param mark formData int true "Updated mark"
// @Param comment formData string false "Updated comment"
// @Success 200 {object} models.Rating "Rating updated successfully"
// @Failure 400 {string} string "Missing or invalid fields in the request"
// @Failure 403 {string} string "User is not authorized to modify this rating"
// @Failure 404 {string} string "Rating not found"
// @Failure 500 {string} string "Internal server error"
// @Router /ratings/{id} [put]
func (h *RatingHandler) UpdateRatingHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	ratingID := chi.URLParam(r, "id")
	markStr := r.FormValue("mark")
	comment := r.FormValue("comment")

	if markStr == "" {
		http.Error(w, "Mark is required", http.StatusBadRequest)
		return
	}

	mark, err := strconv.ParseInt(markStr, 10, 8)
	if err != nil {
		http.Error(w, "Invalid mark value", http.StatusBadRequest)
		return
	}

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "error getting claims", http.StatusInternalServerError)
		return
	}
	userID := claims["id"].(uint)

	existingRating, err := h.ratingQueries.FindRatingByID(ratingID)
	if err != nil {
		http.Error(w, "error finding rating", http.StatusNotFound)
		return
	}

	if existingRating.UserID != userID {
		http.Error(w, "user is not authorized to modify this rating", http.StatusForbidden)
		return
	}

	updatedRating, err := h.ratingQueries.UpdateRating(ratingID, int8(mark), comment)
	if err != nil {
		http.Error(w, "Error updating rating", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedRating)
}

// FetchAllRatingsHandler retrieves all ratings from the database.
// @Summary Fetch all ratings
// @Description Retrieve all ratings from the database
// @Tags ratings
// @Produce json
// @Success 200 {array} models.Rating "List of all ratings"
// @Failure 500 {string} string "Internal server error"
// @Router /ratings [get]
func (h *RatingHandler) FetchAllRatingsHandler(w http.ResponseWriter, r *http.Request) {
	ratings, err := h.ratingQueries.GetAllRatings()
	if err != nil {
		http.Error(w, "error getting ratings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ratings)
}

// DeleteRatingHandler handles the deletion of a rating.
// @Summary Delete rating
// @Description Delete an existing rating
// @Tags ratings
// @Param id path string true "ID of the rating to delete"
// @Success 204 {string} string "Rating deleted successfully"
// @Failure 403 {string} string "User is not authorized to delete this rating"
// @Failure 404 {string} string "Rating not found"
// @Failure 500 {string} string "Internal server error"
// @Router /ratings/{id} [delete]
func (h *RatingHandler) DeleteRatingHandler(w http.ResponseWriter, r *http.Request) {
	ratingID := chi.URLParam(r, "id")

	rating, err := h.ratingQueries.FindRatingByID(ratingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Rating not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error finding rating", http.StatusInternalServerError)
		return
	}

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "Error getting claims", http.StatusInternalServerError)
		return
	}
	userID := claims["id"].(uint)

	if rating.UserID != userID {
		http.Error(w, "User is not authorized to delete this rating", http.StatusForbidden)
		return
	}

	if err := h.ratingQueries.DeleteRating(ratingID); err != nil {
		http.Error(w, "Error deleting rating", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
