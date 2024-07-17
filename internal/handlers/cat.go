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
	"strings"
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
// @Accept json
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
// @Param RaceID formData string true "RaceID"
// @Param Description formData string false "Description"
// @Param Reserved formData string true "Reserved"
// @Param UserID formData string true "User ID"
// @Param PublishedAs formData string true "Published As" // New parameter
// @Param uploaded_file formData file true "Image"
// @Success 201 {object} models.Cats "cat created successfully"
// @Failure 400 {string} string "all fields are required"
// @Failure 500 {string} string "error creating cat"
// @Router /cats [post]
func (h *CatHandler) CatCreationHandler(w http.ResponseWriter, r *http.Request) {
	var (
		name, birthDateStr, sexe, lastVaccineStr, lastVaccineName, color, behavior, sterilizedStr, race, description, reservedStr, userID, publishedAs string
		pictures                                                                                                                                       []string
	)

	contentType := r.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/json") {
		var requestData map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		name, _ = requestData["name"].(string)
		birthDateStr, _ = requestData["BirthDate"].(string)
		sexe, _ = requestData["sexe"].(string)
		lastVaccineStr, _ = requestData["LastVaccine"].(string)
		lastVaccineName, _ = requestData["LastVaccineName"].(string)
		color, _ = requestData["Color"].(string)
		behavior, _ = requestData["Behavior"].(string)
		sterilizedStr, _ = requestData["Sterilized"].(string)
		race, _ = requestData["RaceID"].(string)
		description, _ = requestData["Description"].(string)
		reservedStr, _ = requestData["Reserved"].(string)
		userID, _ = requestData["UserID"].(string)
		publishedAs, _ = requestData["PublishedAs"].(string) // New field
		if uploadedFiles, ok := requestData["uploaded_file"].([]interface{}); ok {
			pictures = convertInterfaceSliceToStringSlice(uploadedFiles)
		}
	} else {
		err := r.ParseMultipartForm(10 << 20) // 10 MB is the max memory size
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		name = r.FormValue("name")
		birthDateStr = r.FormValue("BirthDate")
		sexe = r.FormValue("sexe")
		lastVaccineStr = r.FormValue("LastVaccine")
		lastVaccineName = r.FormValue("LastVaccineName")
		color = r.FormValue("Color")
		behavior = r.FormValue("Behavior")
		sterilizedStr = r.FormValue("Sterilized")
		race = r.FormValue("RaceID")
		description = r.FormValue("Description")
		reservedStr = r.FormValue("Reserved")
		userID = r.FormValue("UserID")
		publishedAs = r.FormValue("PublishedAs") // New field
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

			pictures = append(pictures, FileURL)
		}
	}

	// Validation des champs obligatoires
	if name == "" || birthDateStr == "" || sexe == "" || color == "" || behavior == "" || sterilizedStr == "" || race == "" || reservedStr == "" || userID == "" || publishedAs == "" {
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

	reserved, err := strconv.ParseBool(reservedStr)
	if err != nil {
		http.Error(w, "invalid Reserved format", http.StatusBadRequest)
		return
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
		PicturesURL:     pictures,
		RaceID:          race,
		Description:     &description,
		Reserved:        reserved,
		UserID:          userID,
		PublishedAs:     publishedAs, // New field
	}

	_, err = h.catQueries.CreateCat(cat)
	if err != nil {
		http.Error(w, "error creating cat", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
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
// @Accept json
// @Accept multipart/form-data
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
// @Param RaceID formData string false "RaceID"
// @Param Description formData string false "Description"
// @Param Reserved formData string false "Reserved"
// @Param UserID formData string false "User ID"
// @Param PublishedAs formData string false "Published As" // New parameter
// @Param uploaded_file formData file false "Image"
// @Success 200 {object} models.Cats "Cat updated successfully"
// @Failure 400 {string} string "Missing or invalid fields in the request"
// @Failure 500 {string} string "Internal server error"
// @Router /cats/{id} [put]
func (h *CatHandler) UpdateCatHandler(w http.ResponseWriter, r *http.Request) {
	var (
		name, birthDateStr, sexe, lastVaccineStr, lastVaccineName, color, behavior, sterilizedStr, race, description, reservedStr, userID, publishedAs string
		pictures                                                                                                                                       []string
	)

	contentType := r.Header.Get("Content-Type")
	catID := chi.URLParam(r, "id")

	if catID == "" {
		http.Error(w, "cat ID is required", http.StatusBadRequest)
		return
	}

	if strings.Contains(contentType, "application/json") {
		var requestData map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		name, _ = requestData["name"].(string)
		birthDateStr, _ = requestData["BirthDate"].(string)
		sexe, _ = requestData["sexe"].(string)
		lastVaccineStr, _ = requestData["LastVaccine"].(string)
		lastVaccineName, _ = requestData["LastVaccineName"].(string)
		color, _ = requestData["Color"].(string)
		behavior, _ = requestData["Behavior"].(string)
		sterilizedStr, _ = requestData["Sterilized"].(string)
		race, _ = requestData["RaceID"].(string)
		description, _ = requestData["Description"].(string)
		reservedStr, _ = requestData["Reserved"].(string)
		userID, _ = requestData["UserID"].(string)
		publishedAs, _ = requestData["PublishedAs"].(string) // New field
		if uploadedFiles, ok := requestData["uploaded_file"].([]interface{}); ok {
			pictures = convertInterfaceSliceToStringSlice(uploadedFiles)
		}
	} else {
		err := r.ParseMultipartForm(10 << 20) // 10 MB is the max memory size
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		name = r.FormValue("name")
		birthDateStr = r.FormValue("BirthDate")
		sexe = r.FormValue("sexe")
		lastVaccineStr = r.FormValue("LastVaccine")
		lastVaccineName = r.FormValue("LastVaccineName")
		color = r.FormValue("Color")
		behavior = r.FormValue("Behavior")
		sterilizedStr = r.FormValue("Sterilized")
		race = r.FormValue("RaceID")
		description = r.FormValue("Description")
		reservedStr = r.FormValue("Reserved")
		userID = r.FormValue("UserID")
		publishedAs = r.FormValue("PublishedAs") // New field
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

			pictures = append(pictures, FileURL)
		}
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

	reserved, err := strconv.ParseBool(reservedStr)
	if err != nil {
		http.Error(w, "invalid Reserved format", http.StatusBadRequest)
		return
	}

	cat, err := h.catQueries.FindCatByID(catID)
	if err != nil {
		http.Error(w, "cat not found", http.StatusNotFound)
		return
	}

	if name != "" {
		cat.Name = name
	}
	if birthDate != nil {
		cat.BirthDate = birthDate
	}
	if sexe != "" {
		cat.Sexe = sexe
	}
	if lastVaccine != nil {
		cat.LastVaccine = lastVaccine
	}
	if lastVaccineName != "" {
		cat.LastVaccineName = lastVaccineName
	}
	if color != "" {
		cat.Color = color
	}
	if behavior != "" {
		cat.Behavior = behavior
	}
	if sterilizedStr != "" {
		cat.Sterilized = sterilized
	}
	if race != "" {
		cat.RaceID = race
	}
	if description != "" {
		cat.Description = &description
	}
	if reservedStr != "" {
		cat.Reserved = reserved
	}
	if userID != "" {
		cat.UserID = userID
	}
	if publishedAs != "" {
		cat.PublishedAs = publishedAs // New field
	}
	if len(pictures) > 0 {
		cat.PicturesURL = pictures
	}

	err = h.catQueries.UpdateCat(cat)
	if err != nil {
		http.Error(w, "error updating cat", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
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

	for _, cat := range cats {
		currentRace, err := h.catQueries.FindRaceByID(cat.RaceID)
		if err != nil {
			http.Error(w, "error fetching the race of a cat", http.StatusInternalServerError)
			return
		}
		cat.RaceID = currentRace.RaceName
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
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

	currentRace, err := h.catQueries.FindRaceByID(cat.RaceID)
	if err != nil {
		http.Error(w, "error fetching the race of a cat", http.StatusInternalServerError)
		return
	}
	cat.RaceID = currentRace.RaceName

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(cat); err != nil {
		http.Error(w, "error encoding cat to JSON", http.StatusInternalServerError)
	}
}

// DeleteCatHandler godoc
// @Summary Delete cat by ID
// @Description Delete a cat by its ID
// @Tags cats
// @Param id path string true "Cat ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Cat ID is required"
// @Failure 404 {string} string "Cat not found"
// @Failure 500 {string} string "Error deleting cat"
// @Router /cats/{id} [delete]
func (h *CatHandler) DeleteCatHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

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

// FindCatsByFilterHandler godoc
// @Summary Get cats by filters
// @Description Retrieve cats using their sex, age or race
// @Tags cats
// @Param raceId query string false "RaceID"
// @Param age query int false "Age"
// @Param sexe query string false "Sexe"
// @Produce  json
// @Success 200 {object} []models.Cats "Found cats"
// @Failure 400 {string} string "An error has occured"
// @Failure 404 {string} string "No cats were found"
// @Failure 500 {string} string "error fetching cats"
// @Router /cats/ [get]
func (h *CatHandler) FindCatsByFilterHandler(w http.ResponseWriter, r *http.Request) {
	var data []*models.Annonce
	params := r.URL.Query()

	raceId := params.Get("raceId")
	sexe := params.Get("sexe")
	age, _ := strconv.Atoi(params.Get("age"))
	assoId, _ := strconv.Atoi(params.Get("assoID"))

	var assoName string

	if assoId != 0 {
		asso, _ := h.catQueries.FindAssociationById(assoId)
		assoName = asso.Name
	}

	cats, err := h.catQueries.GetCatByFilters(raceId, age, sexe, assoName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, fmt.Sprintf("Error in parameters"), http.StatusNotFound)
			return
		}
		http.Error(w, "error fetching cat", http.StatusInternalServerError)
		return
	}

	for _, cat := range cats {
		annonce, fail := h.catQueries.FindAnnonceByCatID(fmt.Sprintf("%d", cat.ID))
		if fail != nil {
			continue
		}

		data = append(data, annonce)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(data)

	if len(data) == 0 {
		http.Error(w, "No cats were found with using the filters.", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "error encoding cat to JSON", http.StatusInternalServerError)
		return
	}
}

// GetCatsByUserHandler godoc
// @Summary Get cats by user ID
// @Description Retrieve all cats for a specific user
// @Tags cats
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {array} models.Cats "List of user's cats"
// @Failure 400 {string} string "User ID is required"
// @Failure 500 {string} string "Error fetching cats"
// @Router /cats/user/{userID} [get]
func (h *CatHandler) GetCatsByUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	cats, err := h.catQueries.FindCatsByUserID(userID)
	if err != nil {
		http.Error(w, "error fetching cats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(cats)
	if err != nil {
		http.Error(w, "error encoding cats to JSON", http.StatusInternalServerError)
		return
	}
}

func convertInterfaceSliceToStringSlice(input []interface{}) []string {
	var output []string
	for _, v := range input {
		if str, ok := v.(string); ok {
			output = append(output, str)
		}
	}
	return output
}
