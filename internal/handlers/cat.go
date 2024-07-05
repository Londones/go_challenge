package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"go-challenge/internal/api"
	"go-challenge/internal/database/queries"
	"go-challenge/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/jinzhu/gorm"
	"github.com/uploadcare/uploadcare-go/ucare"
)

type CatHandler struct {
	catQueries       *queries.DatabaseService
	uploadcareClient ucare.Client
}

func NewCatHandler(catQueries *queries.DatabaseService, uploadcareClient ucare.Client) *CatHandler {
	return &CatHandler{catQueries: catQueries, uploadcareClient: uploadcareClient}
}

// CatCreationHandler godoc
// @Summary Create cat
// @Description Create a new cat with the provided details
// @Tags cats
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Name"
// @Param BirthDate formData string true "Birth Date"
// @Param sexe formData string true "Sexe"
// @Param LastVaccine formData string false "Last Vaccine Date"
// @Param LastVaccineName formData string false "Last Vaccine Name"
// @Param Color formData string true "Color"
// @Param Behavior formData string true "Behavior"
// @Param Sterilized formData string true "Sterilized"
// @Param Race formData string true "Race"
// @Param Description formData string false "Description"
// @Param Reserved formData string true "Reserved"
// @Param AnnonceID formData string true "Annonce ID"
// @Param UserID formData string true "User ID"
// @Param uploaded_file formData file true "Image"
// @Success 201 {object} models.Cats "cat created successfully"
// @Failure 400 {string} string "all fields are required"
// @Failure 500 {string} string "error creating cat"
// @Router /cats [post]
func (h *CatHandler) CatCreationHandler(w http.ResponseWriter, r *http.Request) {
	fileURLs := make([]string, 0)

	err := r.ParseMultipartForm(10 << 20) // 10 MB is the max memory size
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")
	birthDateStr := r.FormValue("BirthDate")
	sexe := r.FormValue("sexe")
	lastVaccineStr := r.FormValue("LastVaccine")
	lastVaccineName := r.FormValue("LastVaccineName")
	color := r.FormValue("Color")
	behavior := r.FormValue("Behavior")
	sterilizedStr := r.FormValue("Sterilized")
	race := r.FormValue("Race")
	description := r.FormValue("Description")
	ReservedStr := r.FormValue("Reserved")
	annonceID := r.FormValue("AnnonceID")
	userID := r.FormValue("UserID")

	if name == "" || birthDateStr == "" || sexe == "" || color == "" || behavior == "" || sterilizedStr == "" || race == "" || ReservedStr == "" || annonceID == "" || userID == "" {
		http.Error(w, "all fields are required", http.StatusBadRequest)
		return
	}

	layout := "02-01-2006"
	var birthDate, lastVaccine *time.Time
	if birthDateStr != "" {
		parsedBirthDate, err := time.Parse(layout, birthDateStr)
		if err != nil {
			http.Error(w, "invalid BirthDate format", http.StatusBadRequest)
			return
		}
		birthDate = &parsedBirthDate
	}

	if lastVaccineStr != "" {
		parsedLastVaccine, err := time.Parse(layout, lastVaccineStr)
		if err != nil {
			http.Error(w, "invalid LastVaccine format", http.StatusBadRequest)
			return
		}
		lastVaccine = &parsedLastVaccine
	}

	sterilized, err := strconv.ParseBool(sterilizedStr)
	if err != nil {
		http.Error(w, "invalid Sterilized format", http.StatusBadRequest)
		return
	}

	Reserved, err := strconv.ParseBool(ReservedStr)
	if err != nil {
		http.Error(w, "invalid Reserved format", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["uploaded_file"]
	for _, header := range files {
		file, err := header.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		ext := filepath.Ext(header.Filename)

		tempFile, err := os.CreateTemp(os.TempDir(), "upload-*"+ext)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer os.Remove(tempFile.Name()) // clean up

		_, err = io.Copy(tempFile, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		FileURL, _, err := api.UploadImage(h.uploadcareClient, tempFile.Name())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fileURLs = append(fileURLs, FileURL)
	}

	cat := &models.Cats{
		Name:            name,
		BirthDate:       birthDate,
		Sexe:            sexe,
		LastVaccine:     lastVaccine,
		LastVaccineName: lastVaccineName,
		Color:           color,
		Behavior:        behavior,
		Sterilized:      sterilized,
		PicturesURL:     fileURLs,
		RaceID:          race,
		Description:     &description,
		Reserved:        Reserved,
		AnnonceID:       annonceID,
		UserID:          userID,
	}

	_, err = h.catQueries.CreateCat(cat)
	if err != nil {
		http.Error(w, "error creating cat", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(cat)
	if err != nil {
		http.Error(w, "error encoding cat to JSON", http.StatusInternalServerError)
		return
	}
}

// UpdateCatHandler godoc
// @Summary Update cat
// @Description Update the details of an existing cat
// @Tags cats
// @Accept x-www-form-urlencoded
// @Produce json
// @Param id path string true "Cat ID"
// @Param name formData string false "Name"
// @Param BirthDate formData string false "Birth Date"
// @Param sexe formData string false "Sexe"
// @Param LastVaccine formData string false "Last Vaccine Date"
// @Param LastVaccineName formData string false "Last Vaccine Name"
// @Param Color formData string false "Color"
// @Param Behavior formData string false "Behavior"
// @Param Sterilized formData string false "Sterilized"
// @Param Race formData string false "Race"
// @Param Description formData string false "Description"
// @Param Reserved formData string false "Reserved"
// @Param AnnonceID formData string false "Annonce ID"
// @Success 200 {object} models.Cats "Cat updated successfully"
// @Failure 400 {string} string "Missing or invalid fields in the request"
// @Failure 500 {string} string "Internal server error"
// @Router /cats/{id} [put]
func (h *CatHandler) UpdateCatHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	catID := chi.URLParam(r, "id")
	if catID == "" {
		http.Error(w, "cat ID is required", http.StatusBadRequest)
		return
	}

	cat, err := h.catQueries.FindCatByID(catID)
	if err != nil {
		http.Error(w, "cat not found", http.StatusNotFound)
		return
	}

	if name := r.FormValue("name"); name != "" {
		cat.Name = name
	}
	if birthDateStr := r.FormValue("BirthDate"); birthDateStr != "" {
		layout := "2006-01-02"
		parsedBirthDate, err := time.Parse(layout, birthDateStr)
		if err != nil {
			http.Error(w, "invalid birthDate format", http.StatusBadRequest)
			return
		}
		cat.BirthDate = &parsedBirthDate
	}
	if sexe := r.FormValue("sexe"); sexe != "" {
		cat.Sexe = sexe
	}
	if lastVaccineStr := r.FormValue("LastVaccine"); lastVaccineStr != "" {
		layout := "2006-01-02"
		parsedLastVaccine, err := time.Parse(layout, lastVaccineStr)
		if err != nil {
			http.Error(w, "invalid lastVaccine format", http.StatusBadRequest)
			return
		}
		cat.LastVaccine = &parsedLastVaccine
	}
	if lastVaccineName := r.FormValue("LastVaccineName"); lastVaccineName != "" {
		cat.LastVaccineName = lastVaccineName
	}
	if color := r.FormValue("Color"); color != "" {
		cat.Color = color
	}
	if behavior := r.FormValue("Behavior"); behavior != "" {
		cat.Behavior = behavior
	}
	if sterilizedStr := r.FormValue("Sterilized"); sterilizedStr != "" {
		sterilized, err := strconv.ParseBool(sterilizedStr)
		if err != nil {
			http.Error(w, "invalid sterilized format", http.StatusBadRequest)
			return
		}
		cat.Sterilized = sterilized
	}
	if race := r.FormValue("Race"); race != "" {
		cat.RaceID = race
	}
	if description := r.FormValue("Description"); description != "" {
		cat.Description = &description
	}
	if reservedStr := r.FormValue("Reserved"); reservedStr != "" {
		reserved, err := strconv.ParseBool(reservedStr)
		if err != nil {
			http.Error(w, "invalid reserved format", http.StatusBadRequest)
			return
		}
		cat.Reserved = reserved
	}
	if annonceID := r.FormValue("AnnonceID"); annonceID != "" {
		cat.AnnonceID = annonceID
	}

	err = h.catQueries.UpdateCat(cat)
	if err != nil {
		http.Error(w, "error updating cat", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cat); err != nil {
		http.Error(w, "error encoding cat to JSON", http.StatusInternalServerError)
	}
}

// GetAllCatsHandler godoc
// @Summary Get all cats
// @Description Retrieve a list of all cats
// @Tags cats
// @Produce json
// @Success 200 {array} models.Cats "List of cats"
// @Failure 500 {string} string "error fetching cats"
// @Router /cats [get]
func (h *CatHandler) GetAllCatsHandler(w http.ResponseWriter, r *http.Request) {
	cats, err := h.catQueries.GetAllCats()
	if err != nil {
		http.Error(w, "error fetching cats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cats)
	if err != nil {
		http.Error(w, "error encoding cats to JSON", http.StatusInternalServerError)
		return
	}
}

// GetCatByIDHandler godoc
// @Summary Get cat by ID
// @Description Retrieve a cat by its ID
// @Tags cats
// @Produce json
// @Param id path string true "Cat ID"
// @Success 200 {object} models.Cats "Found cat"
// @Failure 400 {string} string "Cat ID is required"
// @Failure 404 {string} string "Cat not found"
// @Failure 500 {string} string "Error fetching cat"
// @Router /cats/{id} [get]
func (h *CatHandler) GetCatByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID of the cat is required", http.StatusBadRequest)
		return
	}

	cat, err := h.catQueries.FindCatByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, fmt.Sprintf("cat with ID %s not found", id), http.StatusNotFound)
			return
		}
		http.Error(w, "error fetching cat", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cat); err != nil {
		http.Error(w, "error encoding cat to JSON", http.StatusInternalServerError)
	}
}

// DeleteCatHandler godoc
// @Summary Delete cat by ID
// @Description Delete a cat by its ID
// @Tags cats
// @Param id query string true "Cat ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Cat ID is required"
// @Failure 404 {string} string "Cat not found"
// @Failure 500 {string} string "Error deleting cat"
// @Router /cats/{id} [delete]
func (h *CatHandler) DeleteCatHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get("id")

	if id == "" {
		http.Error(w, "cat ID is required", http.StatusBadRequest)
		return
	}

	err := h.catQueries.DeleteCatByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, fmt.Sprintf("cat with ID %s not found", id), http.StatusNotFound)
			return
		}
		http.Error(w, "error deleting cat", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CatHandler) FindCatsByFilterHandler(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()

	raceId := params.Get("raceId")
	sexe := params.Get("sexe")
	age, _ := strconv.Atoi(params.Get("age"))

	cats, err := h.catQueries.GetCatByFilters(raceId, age, sexe)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, fmt.Sprintf("Error in parameters"), http.StatusNotFound)
			return
		}
		http.Error(w, "error fetching cat", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cats)
	if err != nil {
		http.Error(w, "error encoding cat to JSON", http.StatusInternalServerError)
		return
	}
}
