package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jinzhu/gorm"
	"github.com/uploadcare/uploadcare-go/ucare"
	"go-challenge/internal/database/queries"
	"go-challenge/internal/models"
	"net/http"
)

type RaceHandler struct {
	raceQueries      *queries.DatabaseService
	uploadcareClient ucare.Client
}

func NewRaceHandler(raceQueries *queries.DatabaseService, uploadcareClient ucare.Client) *RaceHandler {
	return &RaceHandler{raceQueries: raceQueries, uploadcareClient: uploadcareClient}
}

// GetAllRaceHandler godoc
// @Summary Get all races
// @Description Retrieve a list of all race
// @Tags race
// @Produce  json
// @Success 200 {array} models.Races "List of race"
// @Failure 500 {string} string "error fetching races"
// @Router /races [get]
func (h *RaceHandler) GetAllRaceHandler(w http.ResponseWriter, r *http.Request) {

	races, err := h.raceQueries.GetAllRace()
	if err != nil {
		http.Error(w, "error fetching races", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(races)
	if err != nil {
		http.Error(w, "error encoding races to JSON", http.StatusInternalServerError)
		return
	}
}

// GetRaceByIDHandler godoc
// @Summary Get a specific race using its id
// @Description Retrieve one specific race
// @Tags race
// @Produce json
// @Param id path string true "ID of the race to retrieve"
// @Success 200 {object} models.Races "Race detail"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 404 {string} string "Race not found"
// @Failure 500 {string} string "Internal server error"
// @Router /race/{id} [get]
func (h *RaceHandler) GetRaceByIDHandler(w http.ResponseWriter, r *http.Request) {

	raceID := chi.URLParam(r, "id")
	if raceID == "" {
		http.Error(w, "ID of the race is required", http.StatusBadRequest)
		return
	}

	race, err := h.raceQueries.FindRaceByID(raceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Race not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving race", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(race)
}

// UpdateRaceHandler godoc
// @Summary Search then update a race using its ID
// @Description Update a race name
// @Tags race
// @Produce json
// @Accept x-www-form-urlencoded
// @Param id path string true "ID of the race to update"
// @Param raceName query string true "New race name"
// @Security ApiKeyAuth
// @Success 200 {object} models.Races "Race updated successfully"
// @Failure 400 {string} string "Missing or invalid fields in the request"
// @Failure 403 {string} string "User is not authorized to update this race"
// @Failure 404 {string} string "Race not found"
// @Failure 500 {string} string "Internal server error"
// @Router /race/{id} [put]
func (h *RaceHandler) UpdateRaceHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	raceID := chi.URLParam(r, "id")
	if raceID == "" {
		http.Error(w, "race ID is required", http.StatusBadRequest)
		return
	}

	race, err := h.raceQueries.FindRaceByID(raceID)
	if err != nil {
		http.Error(w, "race not found", http.StatusNotFound)
		return
	}

	if name := r.FormValue("raceName"); name != "" {
		race.RaceName = name
	}

	err = h.raceQueries.UpdateRace(race)
	if err != nil {
		http.Error(w, "error updating race", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(race); err != nil {
		http.Error(w, "error encoding race to JSON", http.StatusInternalServerError)
	}
}

// RaceCreationHandler godoc
// @Summary Create a new Race
// @Description Create a new Race from a Form
// @Tags race
// @Accept  x-www-form-urlencoded
// @Produce  json
// @Param name formData string true "Name"
// @Success 201 {object} models.Races "Race created	successfully"
// @Failure 400 {string} string "all fields are required"
// @Failure 500 {string} string "error creating race"
// @Router /race [post]
func (h *RaceHandler) RaceCreationHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")

	if name == "" {
		http.Error(w, "all fields are required", http.StatusBadRequest)
		return
	}

	race := &models.Races{
		RaceName: name,
	}

	_, err = h.raceQueries.CreateRace(race)
	if err != nil {
		http.Error(w, "error creating race", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(race)
	if err != nil {
		http.Error(w, "error encoding race to JSON", http.StatusInternalServerError)
		return
	}
}

// DeleteRaceHandler godoc
// @Summary Delete a race
// @Description Delete a race using its ID
// @Tags race
// @Param id query string true "ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "race ID is required"
// @Failure 404 {string} string "race not found"
// @Failure 500 {string} string "error deleting race"
// @Router /race/{id} [delete]
func (h *RaceHandler) DeleteRaceHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	id := params.Get("id")

	if id == "" {
		http.Error(w, "race ID is required", http.StatusBadRequest)
		return
	}

	err := h.raceQueries.DeleteRace(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, fmt.Sprintf("race with ID %s not found", id), http.StatusNotFound)
			return
		}
		http.Error(w, "error deleting race", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
