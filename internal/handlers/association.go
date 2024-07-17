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
	"strings"

	"go-challenge/internal/api"
	"go-challenge/internal/database/queries"
	"go-challenge/internal/models"
	"go-challenge/internal/config"

	"github.com/go-chi/chi/v5"
	// "github.com/gorilla/schema"
	"github.com/lib/pq" // Ajout de l'importation du package pq
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
// @Accept json
// @Produce json
// @Param name formData string true "Name"
// @Param addressRue formData string true "AddressRue"
// @Param cp formData string true "CP"
// @Param ville formData string true "Ville"
// @Param phone formData string true "Phone"
// @Param email formData string true "Email"
// @Param ownerId formData string true "OwnerID"
// @Param members formData string false "Comma-separated list of members IDs"
// @Param kbisFile formData file true "PDF file"
// @Success 201 {object} models.Association "Successfully created association"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /associations [post]
func (h *AssociationHandler) CreateAssociationHandler(w http.ResponseWriter, r *http.Request) {
	var association models.Association
	var members []string

	fmt.Println("Starting CreateAssociationHandler")

	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		fmt.Println("Processing JSON payload")
		err := json.NewDecoder(r.Body).Decode(&association)
		if err != nil {
			fmt.Printf("Error decoding JSON: %v\n", err)
			http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		members = association.Members
	} else if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		fmt.Println("Processing multipart/form-data payload")
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			fmt.Printf("Error parsing multipart form: %v\n", err)
			http.Error(w, "Error parsing multipart form: "+err.Error(), http.StatusInternalServerError)
			return
		}

		formData := r.PostForm
		association.Name = formData.Get("name")
		association.AddressRue = formData.Get("addressRue")
		association.Cp = formData.Get("cp")
		association.Ville = formData.Get("ville")
		association.Phone = formData.Get("phone")
		association.Email = formData.Get("email")
		association.OwnerID = formData.Get("ownerId")

		fmt.Printf("Received data: name=%s, addressRue=%s, cp=%s, ville=%s, phone=%s, email=%s, ownerId=%s\n",
			association.Name, association.AddressRue, association.Cp, association.Ville, association.Phone, association.Email, association.OwnerID)

		file, handler, err := r.FormFile("kbisFile")
		if err != nil {
			fmt.Printf("Error retrieving file: %v\n", err)
			http.Error(w, "Error retrieving file: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		if handler.Header.Get("Content-Type") != "application/pdf" {
			fmt.Println("Invalid content type for kbisFile, expected application/pdf")
			http.Error(w, "Invalid content type for kbisFile, expected application/pdf", http.StatusBadRequest)
			return
		}

		ext := filepath.Ext(handler.Filename)
		tempFile, err := os.CreateTemp(os.TempDir(), "upload-*"+ext)
		if err != nil {
			fmt.Printf("Error creating temp file: %v\n", err)
			http.Error(w, "Error creating temp file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer os.Remove(tempFile.Name())

		_, err = io.Copy(tempFile, file)
		if err != nil {
			fmt.Printf("Error copying file: %v\n", err)
			http.Error(w, "Error copying file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		FileURL, _, err := api.UploadFilePDF(h.uploadcareClient, tempFile.Name())
		if err != nil {
			fmt.Printf("Error uploading file: %v\n", err)
			http.Error(w, "Error uploading file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		association.KbisFile = FileURL

		membersStr := formData.Get("members")
		if membersStr != "" {
			members = strings.Split(membersStr, ",")
		}
		association.Members = pq.StringArray(members)
	} else {
		fmt.Println("Unsupported content type")
		http.Error(w, "Unsupported content type", http.StatusBadRequest)
		return
	}

	verified := false
	association.Verified = &verified

	if err := h.associationQueries.CreateAssociation(&association); err != nil {
		fmt.Printf("Error creating association: %v\n", err)
		http.Error(w, "Error creating association: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(association)

	fmt.Println("Association created successfully")
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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(associations)
}

// @Summary Get associations by user ID
// @Description Retrieve all associations for a specific user where the user is either the owner or a member
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

	owner, err := h.associationQueries.FindUserByAssociationID(associationID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	notificationToken, err := h.associationQueries.GetNotificationTokenByUserID(owner.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	payload := make(map[string]string)
	payload["AssociationID"] = strconv.Itoa(associationID)
	payload["Verified"] = strconv.FormatBool(*association.Verified)
	SendToToken(config.GetFirebaseApp(), notificationToken.Token, "Votre association a été vérifié", "Vérification Association", payload)

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

// @Summary Update an association
// @Description Update all fields of an association with the given ID
// @Tags associations
// @Accept multipart/form-data
// @Accept json
// @Produce json
// @Param id path int true "Association ID"
// @Param name formData string false "Name"
// @Param addressRue formData string false "AddressRue"
// @Param cp formData string false "CP"
// @Param ville formData string false "Ville"
// @Param phone formData string false "Phone"
// @Param email formData string false "Email"
// @Param kbisFile formData file false "PDF file"
// @Param members formData string false "Comma-separated list of members IDs"
// @Param association body models.Association false "Association payload"
// @Success 200 {object} models.Association "Successfully updated association"
// @Failure 400 {object} string "Bad Request: Invalid association ID or Invalid content type for kbisFile, expected application/pdf"
// @Failure 500 {object} string "Internal Server Error"
// @Router /associations/{id} [put]
func (h *AssociationHandler) UpdateAssociationHandler(w http.ResponseWriter, r *http.Request) {
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

	var association models.Association
	var members []string

	existingAssociation, err := h.associationQueries.FindAssociationById(associationID)
	if err != nil {
		http.Error(w, "Error finding association: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		err := json.NewDecoder(r.Body).Decode(&association)
		if err != nil {
			http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		members = association.Members
	} else if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			http.Error(w, "Error parsing multipart form: "+err.Error(), http.StatusInternalServerError)
			return
		}

		formData := r.PostForm
		association.Name = formData.Get("name")
		association.AddressRue = formData.Get("addressRue")
		association.Cp = formData.Get("cp")
		association.Ville = formData.Get("ville")
		association.Phone = formData.Get("phone")
		association.Email = formData.Get("email")
		association.OwnerID = formData.Get("ownerId")

		file, handler, err := r.FormFile("kbisFile")
		if err == nil {
			defer file.Close()

			if handler.Header.Get("Content-Type") != "application/pdf" {
				http.Error(w, "Invalid content type for kbisFile, expected application/pdf", http.StatusBadRequest)
				return
			}

			ext := filepath.Ext(handler.Filename)
			tempFile, err := os.CreateTemp(os.TempDir(), "upload-*"+ext)
			if err != nil {
				http.Error(w, "Error creating temp file: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer os.Remove(tempFile.Name())

			_, err = io.Copy(tempFile, file)
			if err != nil {
				http.Error(w, "Error copying file: "+err.Error(), http.StatusInternalServerError)
				return
			}

			FileURL, _, err := api.UploadFilePDF(h.uploadcareClient, tempFile.Name())
			if err != nil {
				http.Error(w, "Error uploading file: "+err.Error(), http.StatusInternalServerError)
				return
			}
			association.KbisFile = FileURL
		}

		membersStr := formData.Get("members")
		if membersStr != "" {
			members = strings.Split(membersStr, ",")
		}
	} else {
		http.Error(w, "Unsupported content type", http.StatusBadRequest)
		return
	}

	if association.Name != "" {
		existingAssociation.Name = association.Name
	}
	if association.AddressRue != "" {
		existingAssociation.AddressRue = association.AddressRue
	}
	if association.Cp != "" {
		existingAssociation.Cp = association.Cp
	}
	if association.Ville != "" {
		existingAssociation.Ville = association.Ville
	}
	if association.Phone != "" {
		existingAssociation.Phone = association.Phone
	}
	if association.Email != "" {
		existingAssociation.Email = association.Email
	}
	if association.KbisFile != "" {
		existingAssociation.KbisFile = association.KbisFile
	}
	if len(members) > 0 {
		existingAssociation.Members = pq.StringArray(members)
	} else {
		existingAssociation.Members = nil
	}

	if err := h.associationQueries.UpdateAssociation(existingAssociation); err != nil {
		http.Error(w, "Error updating association: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingAssociation)
}
