package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"go-challenge/internal/database/queries"
	"go-challenge/internal/models"

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

	if _, err := os.Stat("internal/uploads/association/"); os.IsNotExist(err) {
		err = os.MkdirAll("internal/uploads/association/", 0755)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	dst, err := os.Create("internal/uploads/association/" + handler.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	association.KbisFile = "internal/uploads/association/" + handler.Filename
	
	verified := false
	association.Verified = &verified

	if err := h.associationQueries.CreateAssociation(&association); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(association)
}
