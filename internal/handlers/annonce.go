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
	"github.com/jinzhu/gorm"
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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(annonce)
}

// ModifyAnnonceHandler godoc
// @Summary Modify annonce
// @Description Modify the title, description, and cat ID of an existing annonce
// @Tags annonces
// @Accept  x-www-form-urlencoded
// @Accept  application/json
// @Produce json
// @Param id path string true "ID of the annonce to modify"
// @Param title formData string false "New title of the annonce"
// @Param description formData string false "New description of the annonce"
// @Param catID formData string false "New cat ID of the annonce"
// @Param title body string false "New title of the annonce"
// @Param description body string false "New description of the annonce"
// @Param catID body string false "New cat ID of the annonce"
// @Success 200 {object} models.Annonce "Annonce updated successfully"
// @Failure 400 {string} string "Missing or invalid fields in the request"
// @Failure 403 {string} string "User is not authorized to modify this annonce"
// @Failure 404 {string} string "Annonce not found"
// @Failure 500 {string} string "Internal server error"
// @Router /annonces/{id} [put]
func (h *AnnonceHandler) ModifyAnnonceHandler(w http.ResponseWriter, r *http.Request) {
	var title, description, catID string

	contentType := r.Header.Get("Content-Type")
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
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		title = r.FormValue("title")
		description = r.FormValue("description")
		catID = r.FormValue("catID")
	}

	if title == "" && description == "" && catID == "" {
		http.Error(w, "At least one of title, description, or catID is required", http.StatusBadRequest)
		return
	}

	annonceID := chi.URLParam(r, "id")
	if annonceID == "" {
		http.Error(w, "Annonce ID is required", http.StatusBadRequest)
		return
	}

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "Error getting claims", http.StatusInternalServerError)
		return
	}
	userID, ok := claims["id"].(string)
	if !ok {
		http.Error(w, "Invalid user ID in claims", http.StatusInternalServerError)
		return
	}

	existingAnnonce, err := h.annonceQueries.FindAnnonceByID(annonceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Annonce not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error finding annonce", http.StatusInternalServerError)
		}
		return
	}

	if existingAnnonce.UserID != userID {
		http.Error(w, "User is not authorized to modify this annonce", http.StatusForbidden)
		return
	}

	// Update the fields if provided
	if title != "" {
		existingAnnonce.Title = title
	}
	if description != "" {
		existingAnnonce.Description = &description
	}
	if catID != "" {
		existingAnnonce.CatID = catID
	}

	err = h.annonceQueries.UpdateAnnonce(existingAnnonce)
	if err != nil {
		http.Error(w, "Error updating annonce", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingAnnonce)
}

// DeleteAnnonceHandler godoc
// @Summary Delete annonce by ID
// @Description Delete an annonce by its ID
// @Tags annonces
// @Param id path string true "Annonce ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Annonce ID is required"
// @Failure 404 {string} string "Annonce not found"
// @Failure 500 {string} string "Error deleting annonce"
// @Router /annonces/{id} [delete]
func (h *AnnonceHandler) DeleteAnnonceHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "annonce ID is required", http.StatusBadRequest)
		return
	}

	err := h.annonceQueries.DeleteRoomByAnnonceID(id)
	if err != nil {
		http.Error(w, "error deleting room", http.StatusInternalServerError)
		return
	}

	err = h.annonceQueries.DeleteAnnonce(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, fmt.Sprintf("annonce with ID %s not found", id), http.StatusNotFound)
			return
		}
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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(annonce)
}
