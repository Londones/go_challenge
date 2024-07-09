package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"go-challenge/internal/database/queries"
	"go-challenge/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/uploadcare/uploadcare-go/ucare"
)

type AssociationHandler struct {
	associationQueries *queries.DatabaseService
	uploadcareClient ucare.Client
	
}

func NewAssociationHandler(associationQueries *queries.DatabaseService, uploadcareClient ucare.Client) *AssociationHandler {
	return &AssociationHandler{associationQueries: associationQueries, uploadcareClient: uploadcareClient}
}

// @Summary Create a new association
// @Description Create a new association with the provided details and uploaded file
// @ID create-association
// @Accept multipart/form-data
// @Produce json
// @Param association body models.Association true "Association details"
// @Param kbisFile formData file true "KBIS file"
// @Success 201 {object} models.Association "Successfully created association"
// @Failure 400 {object} string "Invalid form data"
// @Failure 500 {object} string "Internal server error"
// @Router /associations [post]
func (h *AssociationHandler) CreateAssociationHandler(w http.ResponseWriter, r *http.Request) {
    var association models.Association

    err := r.ParseMultipartForm(10 << 20) // 10 MB
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    formData := r.PostForm
    decoder := schema.NewDecoder()
    err = decoder.Decode(&association, formData)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    file, handler, err := r.FormFile("kbisFile")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()

    verified := false
    association.Verified = &verified

    if err := h.associationQueries.CreateAssociation(&association); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    dirPath := fmt.Sprintf("internal/uploads/association/%d", association.ID)
    if _, err := os.Stat(dirPath); os.IsNotExist(err) {
        err = os.MkdirAll(dirPath, 0755)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }

    dstPath := fmt.Sprintf("%s/%s", dirPath, handler.Filename)
    dst, err := os.Create(dstPath)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer dst.Close()

    if _, err := io.Copy(dst, file); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    association.KbisFile = dstPath

    if err := h.associationQueries.UpdateAssociation(&association); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(association)
}

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

