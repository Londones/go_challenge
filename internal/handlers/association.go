package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"go-challenge/internal/api"
	"go-challenge/internal/database/queries"
	"go-challenge/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/uploadcare/uploadcare-go/ucare"
)

type AssociationHandler struct {
	associationQueries *queries.DatabaseService
	uploadcareClient   ucare.Client
}

func NewAssociationHandler(associationQueries *queries.DatabaseService, uploadcareClient ucare.Client) *AssociationHandler {
	return &AssociationHandler{associationQueries: associationQueries, uploadcareClient: uploadcareClient}
}

// @Summary Create a new association
// @Description Create a new association with the input payload and a PDF file
// @Tags associations
// @Accept multipart/form-data
// @Produce json
// @Param association body models.Association true "Association payload"
// @Param kbisFile formData file true "PDF file"
// @Success 201 {object} models.Association "Successfully created association"
// @Failure 400 {object} string "Bad Request: Error uploading image 2/3, Invalid content type for kbisFile, expected application/pdf"
// @Failure 500 {object} string "Internal Server Error: Error uploading image 1/4/5/6/7/8"
// @Router /associations [post]
func (h *AssociationHandler) CreateAssociationHandler(w http.ResponseWriter, r *http.Request) {
	var association models.Association

	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error uploading image 1: "+err.Error(), http.StatusInternalServerError)
		return
	}

	formData := r.PostForm
	decoder := schema.NewDecoder()
	err = decoder.Decode(&association, formData)
	if err != nil {
		http.Error(w, "Error uploading image 2: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("kbisFile")
	if err != nil {
		http.Error(w, "Error uploading image 3: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	if handler.Header.Get("Content-Type") != "application/pdf" {
		http.Error(w, "Invalid content type for kbisFile, expected application/pdf", http.StatusBadRequest)
		return
	}

	verified := false
	association.Verified = &verified

	if _, err := h.associationQueries.CreateAssociation(&association); err != nil {
		http.Error(w, "Error uploading image 4: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Upload file to uploadcare
	ext := filepath.Ext(handler.Filename)

	tempFile, err := os.CreateTemp(os.TempDir(), "upload-*"+ext)
	if err != nil {
		http.Error(w, "Error uploading image 5: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Error uploading image 6: "+err.Error(), http.StatusInternalServerError)
		return
	}

	FileURL, _, err := api.UploadFilePDF(h.uploadcareClient, tempFile.Name())
	if err != nil {
		http.Error(w, "Error uploading image 7: "+err.Error(), http.StatusInternalServerError)
		return
	}

	association.KbisFile = FileURL

	if err := h.associationQueries.UpdateAssociation(&association); err != nil {
		http.Error(w, "Error uploading image 8: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(association)
}

// @Summary Get all associations
// @Description Retrieve all associations from the database
// @Tags associations
// @Produce json
// @Success 200 {array} models.Association "Successfully retrieved all associations"
// @Failure 500 {object} string "Internal Server Error"
// @Router /associations [get]
func (h *AssociationHandler) GetAllAssociationsHandler(w http.ResponseWriter, r *http.Request) {
	associations, err := h.associationQueries.GetAllAssociations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(associations)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(associations)
}

// @Summary Get associations by user ID
// @Description Retrieve all associations for a specific user
// @Tags associations
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {array} models.Association "Successfully retrieved associations for user"
// @Failure 400 {object} string "Bad Request: Invalid user ID"
// @Failure 500 {object} string "Internal Server Error"
// @Router /users/{userId}/associations [get]
func (h *AssociationHandler) GetUserAssociationsHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	associations, err := h.associationQueries.FindAssociationsByUserId(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(associations)
}

// @Summary Get association by ID
// @Description Retrieve an association by its ID
// @Tags associations
// @Produce json
// @Param id path int true "Association ID"
// @Success 200 {object} models.Association "Successfully retrieved association"
// @Failure 400 {object} string "Bad Request: Invalid association ID"
// @Failure 404 {object} string "Not Found: Association not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /associations/{id} [get]
func (h *AssociationHandler) GetAssociationByIdHandler(w http.ResponseWriter, r *http.Request) {
	associationIDStr := chi.URLParam(r, "id")
	if associationIDStr == "" {
		http.Error(w, "Missing association ID", http.StatusBadRequest)
		return
	}

	associationID, err := strconv.Atoi(associationIDStr)
	if err != nil {
		http.Error(w, "Invalid association ID", http.StatusBadRequest)
		return
	}

	association, err := h.associationQueries.FindAssociationById(associationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(association)
}

// @Summary Update an association's verify status
// @Description Update the verify status of an association with the given ID
// @Tags associations
// @Accept json
// @Produce json
// @Param id path int true "Association ID"
// @Param verified body bool true "Verify status"
// @Success 200 {object} models.Association "Successfully updated association"
// @Failure 400 {object} string "Bad Request: Missing association ID, Invalid association ID"
// @Failure 500 {object} string "Internal Server Error"
// @Router /associations/{id}/verify [put]
func (h *AssociationHandler) UpdateAssociationVerifyStatusHandler(w http.ResponseWriter, r *http.Request) {
	associationIDStr := chi.URLParam(r, "id")
	if associationIDStr == "" {
		http.Error(w, "Missing association ID", http.StatusBadRequest)
		return
	}

	associationID, err := strconv.Atoi(associationIDStr)
	if err != nil {
		http.Error(w, "Invalid association ID", http.StatusBadRequest)
		return
	}

	var body struct {
		Verified *bool `json:"verified"`
	}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	association, err := h.associationQueries.FindAssociationById(associationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	association.Verified = body.Verified

	if err := h.associationQueries.UpdateAssociation(association); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(association)
}

// @Summary Delete an association
// @Description Delete an association by its ID
// @Tags associations
// @Produce json
// @Param id path int true "Association ID"
// @Success 204 "Successfully deleted association"
// @Failure 400 {object} string "Bad Request: Invalid association ID"
// @Failure 404 {object} string "Not Found: Association not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /associations/{id} [delete]
func (h *AssociationHandler) DeleteAssociationHandler(w http.ResponseWriter, r *http.Request) {
	associationIDStr := chi.URLParam(r, "id")
	if associationIDStr == "" {
		http.Error(w, "Missing association ID", http.StatusBadRequest)
		return
	}

	associationID, err := strconv.Atoi(associationIDStr)
	if err != nil {
		http.Error(w, "Invalid association ID", http.StatusBadRequest)
		return
	}

	err = h.associationQueries.DeleteAssociation(associationID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Association not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
