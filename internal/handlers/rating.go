package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go-challenge/internal/database/queries"
	"go-challenge/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jinzhu/gorm"
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

// CreateRatingHandler gère la création d'une nouvelle évaluation.
// @Summary Créer une évaluation
// @Description Crée une nouvelle évaluation avec les détails fournis
// @Tags Ratings
// @Accept x-www-form-urlencoded
// @Accept application/json
// @Produce json
// @Param mark formData int true "Note de l'évaluation"
// @Param comment formData string false "Commentaire sur l'évaluation"
// @Param userID formData string true "ID de l'utilisateur évalué"
// @Success 201 {object} models.Rating "Évaluation créée avec succès"
// @Failure 400 {string} string "Champs manquants ou invalides dans la requête"
// @Failure 500 {string} string "Erreur interne du serveur"
// @Router /ratings [post]
func (h *RatingHandler) CreateRatingHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	var mark int
	var comment, userID string

	if strings.Contains(contentType, "application/json") {
		var data struct {
			Mark    int    `json:"mark"`
			Comment string `json:"comment"`
			UserID  string `json:"userID"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		mark = data.Mark
		comment = data.Comment
		userID = data.UserID
	} else {
		r.ParseForm()
		markStr := r.FormValue("mark")
		comment = r.FormValue("comment")
		userID = r.FormValue("userID")

		if markStr == "" || userID == "" {
			http.Error(w, "La note et l'ID de l'utilisateur sont requis", http.StatusBadRequest)
			return
		}

		var err error
		mark, err = strconv.Atoi(markStr)
		if err != nil {
			http.Error(w, "Valeur de note invalide", http.StatusBadRequest)
			return
		}
	}

	if mark == 0 || userID == "" {
		http.Error(w, "La note et l'ID de l'utilisateur sont requis", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received Data: Mark=%d, Comment=%s, UserID=%s\n", mark, comment, userID)

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des claims", http.StatusInternalServerError)
		return
	}
	authorID := claims["id"].(string)

	rating := &models.Rating{
		Mark:     int8(mark),
		Comment:  comment,
		UserID:   userID,
		AuthorID: authorID,
	}

	if _, err := h.ratingQueries.CreateRating(rating); err != nil {
		http.Error(w, "Erreur lors de la création de l'évaluation", http.StatusInternalServerError)
		return
	}

	createdRating, err := h.ratingQueries.FindRatingByID(fmt.Sprintf("%d", rating.ID))
	if err != nil {
		http.Error(w, "Erreur lors de la récupération de l'évaluation créée", http.StatusInternalServerError)
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

// UpdateRatingHandler gère la mise à jour d'une évaluation existante.
// @Summary Mettre à jour une évaluation
// @Description Met à jour les détails d'une évaluation existante
// @Tags Ratings
// @Accept x-www-form-urlencoded
// @Accept application/json
// @Produce json
// @Param id path string true "ID de l'évaluation à mettre à jour"
// @Param mark formData int true "Note mise à jour"
// @Param comment formData string false "Commentaire mis à jour"
// @Success 200 {object} models.Rating "Évaluation mise à jour avec succès"
// @Failure 400 {string} string "Champs manquants ou invalides dans la requête"
// @Failure 403 {string} string "L'utilisateur n'est pas autorisé à modifier cette évaluation"
// @Failure 404 {string} string "Évaluation non trouvée"
// @Failure 500 {string} string "Erreur interne du serveur"
// @Router /ratings/{id} [put]
func (h *RatingHandler) UpdateRatingHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	var mark int
	var comment string

	ratingID := chi.URLParam(r, "id")

	if strings.Contains(contentType, "application/json") {
		var data struct {
			Mark    int    `json:"mark"`
			Comment string `json:"comment"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		mark = data.Mark
		comment = data.Comment
	} else {
		r.ParseForm()
		markStr := r.FormValue("mark")
		comment = r.FormValue("comment")

		if markStr == "" {
			http.Error(w, "La note est requise", http.StatusBadRequest)
			return
		}

		var err error
		mark, err = strconv.Atoi(markStr)
		if err != nil {
			http.Error(w, "Valeur de note invalide", http.StatusBadRequest)
			return
		}
	}

	if mark == 0 {
		http.Error(w, "La note est requise", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received Data: Mark=%d, Comment=%s\n", mark, comment)

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des claims", http.StatusInternalServerError)
		return
	}
	authorID := claims["id"].(string)

	existingRating, err := h.ratingQueries.FindRatingByID(ratingID)
	if err != nil {
		http.Error(w, "Erreur lors de la recherche de l'évaluation", http.StatusNotFound)
		return
	}

	if existingRating.AuthorID != authorID {
		http.Error(w, "L'utilisateur n'est pas autorisé à modifier cette évaluation", http.StatusForbidden)
		return
	}

	existingRating.Mark = int8(mark)
	existingRating.Comment = comment

	if err := h.ratingQueries.UpdateRating(existingRating); err != nil {
		http.Error(w, "Erreur lors de la mise à jour de l'évaluation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingRating)
}

// FetchAllRatingsHandler récupère toutes les Ratings de la base de données.
// @Summary Récupérer toutes les Ratings
// @Description Récupère toutes les Ratings de la base de données
// @Tags Ratings
// @Produce json
// @Success 200 {array} models.Rating "Liste de toutes les Ratings"
// @Failure 500 {string} string "Erreur interne du serveur"
// @Router /ratings [get]
func (h *RatingHandler) FetchAllRatingsHandler(w http.ResponseWriter, r *http.Request) {
	ratings, err := h.ratingQueries.GetAllRatings()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des évaluations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ratings)
}

// GetRatingByIDHandler récupère une évaluation par son ID.
// @Summary Récupérer une évaluation par ID
// @Description Récupère une évaluation par son ID
// @Tags Ratings
// @Produce json
// @Param id path string true "ID de l'évaluation à récupérer"
// @Success 200 {object} models.Rating "Évaluation récupérée avec succès"
// @Failure 404 {string} string "Évaluation non trouvée"
// @Failure 500 {string} string "Erreur interne du serveur"
// @Router /ratings/{id} [get]
func (h *RatingHandler) GetRatingByIDHandler(w http.ResponseWriter, r *http.Request) {
	ratingID := chi.URLParam(r, "id")

	rating, err := h.ratingQueries.FindRatingByID(ratingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Évaluation non trouvée", http.StatusNotFound)
			return
		}
		http.Error(w, "Erreur lors de la recherche de l'évaluation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rating)
}

// DeleteRatingHandler gère la suppression d'une évaluation.
// @Summary Supprimer une évaluation
// @Description Supprime une évaluation existante
// @Tags Ratings
// @Param id path string true "ID de l'évaluation à supprimer"
// @Success 204 {string} string "Évaluation supprimée avec succès"
// @Failure 403 {string} string "L'utilisateur n'est pas autorisé à supprimer cette évaluation"
// @Failure 404 {string} string "Évaluation non trouvée"
// @Failure 500 {string} string "Erreur interne du serveur"
// @Router /ratings/{id} [delete]
func (h *RatingHandler) DeleteRatingHandler(w http.ResponseWriter, r *http.Request) {
	ratingID := chi.URLParam(r, "id")

	rating, err := h.ratingQueries.FindRatingByID(ratingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Évaluation non trouvée", http.StatusNotFound)
			return
		}
		http.Error(w, "Erreur lors de la recherche de l'évaluation", http.StatusInternalServerError)
		return
	}

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des claims", http.StatusInternalServerError)
		return
	}
	authorID := claims["id"].(string)

	if rating.AuthorID != authorID {
		http.Error(w, "L'utilisateur n'est pas autorisé à supprimer cette évaluation", http.StatusForbidden)
		return
	}

	if err := h.ratingQueries.DeleteRating(ratingID); err != nil {
		http.Error(w, "Erreur lors de la suppression de l'évaluation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetUserRatingsHandler récupère toutes les Ratings pour un utilisateur spécifique.
// @Summary Récupérer les Ratings d'un utilisateur
// @Description Récupère toutes les Ratings pour un utilisateur spécifique
// @Tags Ratings
// @Produce json
// @Param userID path string true "ID de l'utilisateur"
// @Success 200 {array} models.Rating "Liste des Ratings pour l'utilisateur"
// @Failure 404 {string} string "Utilisateur non trouvé"
// @Failure 500 {string} string "Erreur interne du serveur"
// @Router /ratings/user/{userID} [get]
func (h *RatingHandler) GetUserRatingsHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	ratings, err := h.ratingQueries.GetUserRatings(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Utilisateur non trouvé", http.StatusNotFound)
			return
		}
		http.Error(w, "Erreur lors de la récupération des évaluations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ratings)
}

// GetAuthorsRatingsHandler récupère toutes les Ratings créées par un auteur spécifique.
// @Summary Récupérer les Ratings d'un auteur
// @Description Récupère toutes les Ratings créées par un auteur spécifique
// @Tags Ratings
// @Produce json
// @Param authorID path string true "ID de l'auteur"
// @Success 200 {array} models.Rating "Liste des Ratings par l'auteur"
// @Failure 404 {string} string "Auteur non trouvé"
// @Failure 500 {string} string "Erreur interne du serveur"
// @Router /ratings/author/{authorID} [get]
func (h *RatingHandler) GetAuthorsRatingsHandler(w http.ResponseWriter, r *http.Request) {
	authorID := chi.URLParam(r, "authorID")

	ratings, err := h.ratingQueries.GetAuthorRatings(authorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Auteur non trouvé", http.StatusNotFound)
			return
		}
		http.Error(w, "Erreur lors de la récupération des évaluations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ratings)
}
