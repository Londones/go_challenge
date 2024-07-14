package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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

// UpdateCatHandler godoc
// @Summary Update cat
// @Description Update the details of an existing cat
// @Tags cats
// @Accept x-www-form-urlencoded
// @Accept json
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
// @Param AnnonceID formData string false "Annonce ID"
// @Param body body models.Cats false "Cat object"
// @Success 200 {object} models.Cats "Cat updated successfully"
// @Failure 400 {string} string "Missing or invalid fields in the request"
// @Failure 500 {string} string "Internal server error"
// @Router /cats/{id} [put]
func (h *CatHandler) UpdateCatHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	var updateData models.Cats

	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if name := r.FormValue("name"); name != "" {
			updateData.Name = name
		}
		if birthDateStr := r.FormValue("BirthDate"); birthDateStr != "" {
			layout := "2006-01-02"
			parsedBirthDate, err := time.Parse(layout, birthDateStr)
			if err != nil {
				http.Error(w, "invalid birthDate format", http.StatusBadRequest)
				return
			}
			updateData.BirthDate = &parsedBirthDate
		}
		if sexe := r.FormValue("sexe"); sexe != "" {
			updateData.Sexe = sexe
		}
		if lastVaccineStr := r.FormValue("LastVaccine"); lastVaccineStr != "" {
			layout := "2006-01-02"
			parsedLastVaccine, err := time.Parse(layout, lastVaccineStr)
			if err != nil {
				http.Error(w, "invalid lastVaccine format", http.StatusBadRequest)
				return
			}
			updateData.LastVaccine = &parsedLastVaccine
		}
		if lastVaccineName := r.FormValue("LastVaccineName"); lastVaccineName != "" {
			updateData.LastVaccineName = lastVaccineName
		}
		if color := r.FormValue("Color"); color != "" {
			updateData.Color = color
		}
		if behavior := r.FormValue("Behavior"); behavior != "" {
			updateData.Behavior = behavior
		}
		if sterilizedStr := r.FormValue("Sterilized"); sterilizedStr != "" {
			sterilized, err := strconv.ParseBool(sterilizedStr)
			if err != nil {
				http.Error(w, "invalid sterilized format", http.StatusBadRequest)
				return
			}
			updateData.Sterilized = sterilized
		}
		if raceID := r.FormValue("RaceID"); raceID != "" {
			updateData.RaceID = raceID
		}
		if description := r.FormValue("Description"); description != "" {
			updateData.Description = &description
		}
		if reservedStr := r.FormValue("Reserved"); reservedStr != "" {
			reserved, err := strconv.ParseBool(reservedStr)
			if err != nil {
				http.Error(w, "invalid reserved format", http.StatusBadRequest)
				return
			}
			updateData.Reserved = reserved
		}
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

	if updateData.Name != "" {
		cat.Name = updateData.Name
	}
	if updateData.BirthDate != nil {
		cat.BirthDate = updateData.BirthDate
	}
	if updateData.Sexe != "" {
		cat.Sexe = updateData.Sexe
	}
	if updateData.LastVaccine != nil {
		cat.LastVaccine = updateData.LastVaccine
	}
	if updateData.LastVaccineName != "" {
		cat.LastVaccineName = updateData.LastVaccineName
	}
	if updateData.Color != "" {
		cat.Color = updateData.Color
	}
	if updateData.Behavior != "" {
		cat.Behavior = updateData.Behavior
	}
	if updateData.Sterilized {
		cat.Sterilized = updateData.Sterilized
	}
	if updateData.RaceID != "" {
		cat.RaceID = updateData.RaceID
	}
	if updateData.Description != nil {
		cat.Description = updateData.Description
	}
	if updateData.Reserved {
		cat.Reserved = updateData.Reserved
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

	for _, cat := range cats {
		currentRace, err := h.catQueries.FindRaceByID(cat.RaceID)
		if err != nil {
			http.Error(w, "error fetching the race of a cat", http.StatusInternalServerError)
			return
		}
		cat.RaceID = currentRace.RaceName
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

	currentRace, err := h.catQueries.FindRaceByID(cat.RaceID)
	if err != nil {
		http.Error(w, "error fetching the race of a cat", http.StatusInternalServerError)
		return
	}
	cat.RaceID = currentRace.RaceName

	w.Header().Set("Content-Type", "application/json")
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

	cats, err := h.catQueries.GetCatByFilters(raceId, age, sexe)
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

	w.Header().Set("Content-Type", "application/json")
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

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cats)
	if err != nil {
		http.Error(w, "error encoding cats to JSON", http.StatusInternalServerError)
		return
	}
}
