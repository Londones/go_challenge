package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"go-challenge/internal/database/queries"
	"go-challenge/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

type AnnonceHandler struct {
	annonceQueries *queries.DatabaseService
	userQueries    *queries.DatabaseService
	catQueries     *queries.DatabaseService
}

func NewAnnonceHandler(annonceQueries, userQueries, catQueries *queries.DatabaseService) *AnnonceHandler {
	return &AnnonceHandler{annonceQueries: annonceQueries, userQueries: userQueries, catQueries: catQueries}
}

// AnnonceCreationHandler godoc
// @Summary Create annonces
// @Description Create a new annonce with the provided details
// @Tags annonces
// @Accept  x-www-form-urlencoded
// @Produce  json
// @Param title formData string true "Title of the annonce"
// @Param description formData string true "Description of the annonce"
// @Param catID formData string true "Cat ID"
// @Success 201 {object} models.Annonce "Annonce created successfully"
// @Failure 400 {string} string "Missing or invalid fields in the request"
// @Failure 500 {string} string "Internal server error"
// @Router /annonces [post]
func (h *AnnonceHandler) AnnonceCreationHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	var title, description, catID string

	if strings.Contains(contentType, "application/json") {
		var data struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			CatID       string `json:"catID"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		title = data.Title
		description = data.Description
		catID = data.CatID
	} else {
		r.ParseForm()
		title = r.FormValue("title")
		description = r.FormValue("description")
		catID = r.FormValue("catID")
	}

	if title == "" || description == "" || catID == "" {
		http.Error(w, "title, description, and catID are required", http.StatusBadRequest)
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

	cat, err := h.catQueries.FindCatByID(catID)
	if err != nil {
		http.Error(w, "error finding cat", http.StatusInternalServerError)
		return
	}

	annonce := &models.Annonce{
		Title:       title,
		Description: &description,
		UserID:      user.ID,
		CatID:       fmt.Sprintf("%d", cat.ID),
	}

	annonceID, err := h.annonceQueries.CreateAnnonce(annonce)
	if err != nil {
		http.Error(w, "error creating annonce", http.StatusInternalServerError)
		return
	}

	createdAnnonce, err := h.annonceQueries.FindAnnonceByID(fmt.Sprintf("%d", annonceID))
	if err != nil {
		http.Error(w, "error retrieving created annonce", http.StatusInternalServerError)
		return
	}

	response := struct {
		Success string          `json:"success"`
		Annonce *models.Annonce `json:"annonce"`
	}{
		Success: "true",
		Annonce: createdAnnonce,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetAllAnnoncesHandler godoc
// @Summary Get all annonces
// @Description Retrieve all annonces from the database
// @Tags annonces
// @Produce json
// @Success 200 {array} models.Annonce "List of annonces"
// @Failure 500 {string} string "Internal server error"
// @Router /annonces [get]
func (h *AnnonceHandler) GetAllAnnoncesHandler(w http.ResponseWriter, r *http.Request) {
	annonces, err := h.annonceQueries.GetAllAnnonces()
	if err != nil {
		http.Error(w, "error getting annonces", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(annonces)
}

// GetUserAnnoncesHandler godoc
// @Summary Get user's annonces
// @Description Retrieve all annonces for a specific user from the database
// @Tags annonces
// @Produce json
// @Param id query string true "User ID"
// @Success 200 {array} models.Annonce "List of user's annonces"
// @Failure 400 {string} string "Bad request - missing userID parameter"
// @Failure 500 {string} string "Internal server error"
// @Router /users/annonces/{id} [get]
func (h *AnnonceHandler) GetUserAnnoncesHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	var userID string

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "error getting claims", http.StatusInternalServerError)
		return
	}

	if id := params.Get("id"); id != "" {
		userID = id
	} else if claims["id"] != nil {
		userID = claims["id"].(string)
	}

	if userID == "" {
		http.Error(w, "missing userID parameter", http.StatusBadRequest)
		return
	}

	annonces, err := h.annonceQueries.GetUserAnnonces(userID)
	if err != nil {
		http.Error(w, "error getting user's annonces", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(annonces)
	if err != nil {
		http.Error(w, "error encoding annonces to JSON", http.StatusInternalServerError)
		return
	}
}

// GetAnnonceByIDHandler godoc
// @Summary Get an annonce by ID
// @Description Retrieve an annonce from the database by its ID
// @Tags annonces
// @Produce json
// @Param id path string true "ID of the annonce to retrieve"
// @Success 200 {object} models.Annonce "Annonce details"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 404 {string} string "Annonce not found"
// @Failure 500 {string} string "Internal server error"
// @Router /annonces/{id} [get]
func (h *AnnonceHandler) GetAnnonceByIDHandler(w http.ResponseWriter, r *http.Request) {
	annonceID := chi.URLParam(r, "id")
	if annonceID == "" {
		http.Error(w, "ID of the annonce is required", http.StatusBadRequest)
		return
	}

	print("Annonce ID: ", annonceID)

	annonce, err := h.annonceQueries.FindAnnonceByID(annonceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Annonce not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving annonce", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(annonce)
}

// ModifyDescriptionAnnonceHandler godoc
// @Summary Modify annonce description
// @Description Modify the description of an existing annonce
// @Tags annonces
// @Accept  x-www-form-urlencoded
// @Produce json
// @Param id path string true "ID of the annonce to modify"
// @Param description formData string true "New description of the annonce"
// @Success 200 {object} models.Annonce "Annonce updated successfully"
// @Failure 400 {string} string "Missing or invalid fields in the request"
// @Failure 403 {string} string "User is not authorized to modify this annonce"
// @Failure 404 {string} string "Annonce not found"
// @Failure 500 {string} string "Internal server error"
// @Router /annonces/{id} [put]
func (h *AnnonceHandler) ModifyDescriptionAnnonceHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	annonceID := chi.URLParam(r, "id")
	description := r.FormValue("description")

	if description == "" {
		http.Error(w, "description is required", http.StatusBadRequest)
		return
	}

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "error getting claims", http.StatusInternalServerError)
		return
	}
	userID := claims["id"].(string)

	existingAnnonce, err := h.annonceQueries.FindAnnonceByID(annonceID)
	if err != nil {
		http.Error(w, "error finding annonce", http.StatusNotFound)
		return
	}

	if existingAnnonce.UserID != userID {
		http.Error(w, "user is not authorized to modify this annonce", http.StatusForbidden)
		return
	}

	updatedAnnonce, err := h.annonceQueries.UpdateAnnonceDescription(annonceID, description)
	if err != nil {
		http.Error(w, "Error updating annonce", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedAnnonce)
}

// DeleteAnnonceHandler godoc
// @Summary Delete annonce
// @Description Delete an existing annonce
// @Tags annonces
// @Param id path string true "ID of the annonce to delete"
// @Success 204 {string} string "Annonce deleted successfully"
// @Failure 403 {string} string "User is not authorized to delete this annonce"
// @Failure 404 {string} string "Annonce not found"
// @Failure 500 {string} string "Internal server error"
// @Router /annonces/{id} [delete]
func (h *AnnonceHandler) DeleteAnnonceHandler(w http.ResponseWriter, r *http.Request) {
	annonceID := chi.URLParam(r, "id")

	annonce, err := h.annonceQueries.FindAnnonceByID(annonceID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "annonce not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error finding annonce", http.StatusInternalServerError)
		return
	}

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "error getting claims", http.StatusInternalServerError)
		return
	}
	userID := claims["id"].(string)

	if annonce.UserID != userID {
		http.Error(w, "user is not authorized to modify this annonce", http.StatusForbidden)
		return
	}

	if err := h.annonceQueries.DeleteAnnonce(annonceID); err != nil {
		http.Error(w, "error deleting annonce", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// FetchAnnonceByCatIDHandler godoc
// @Summary Get an annonce by Cat ID
// @Description Retrieve an annonce from the database by its Cat ID
// @Tags annonces
// @Produce json
// @Param catID path string true "Cat ID of the annonce to retrieve"
// @Success 200 {object} models.Annonce "Annonce details"
// @Failure 400 {string} string "Invalid Cat ID format"
// @Failure 404 {string} string "Annonce not found"
// @Failure 500 {string} string "Internal server error"
// @Router /annonces/cats/{catID} [get]
func (h *AnnonceHandler) FetchAnnonceByCatIDHandler(w http.ResponseWriter, r *http.Request) {
	catID := chi.URLParam(r, "catID")
	if catID == "" {
		http.Error(w, "Cat ID is required", http.StatusBadRequest)
		return
	}

	annonce, err := h.annonceQueries.FindAnnonceByCatID(catID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Annonce not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving annonce", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(annonce)
}

// GetAddressFromUserID godoc
// @Summary Get the user address from user ID
// @Description Get the address from the user ID
// @Tags annonces
// @Produce json
// @Param userID path string true "ID of the annonce's user"
// @Success 200 {object} String "Address"
// @Failure 400 {string} string "Invalid annonce ID format"
// @Failure 404 {string} string "Annonce not found"
// @Failure 500 {string} string "Internal server error"
// @Router /annonces/address/{id} [get]
func (h *AnnonceHandler) GetAddressFromUserID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if annonceID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	address, err := h.annonceQueries.GetAddressFromAnnonceID(userID)
}
