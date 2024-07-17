package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"go-challenge/internal/database/queries"
	"go-challenge/internal/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/uploadcare/uploadcare-go/ucare"
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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
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
// @Router /races/{id} [get]
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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(race)
}

// UpdateRaceHandler updates a race
// @Summary Update a race
// @Description Update a race by ID
// @Tags race
// @Accept  json
// @Produce  json
// @Param id path string true "Race ID"
// @Param body body models.Races true "Race object"
// @Success 200 {object} models.Races "Successfully updated race"
// @Failure 400 {object} string "Invalid ID supplied"
// @Failure 404 {object} string "Race not found"
// @Failure 400 {object} string "Invalid JSON body"
// @Failure 500 {object} string "Error updating race"
// @Router /races/{id} [put]
func (h *RaceHandler) UpdateRaceHandler(w http.ResponseWriter, r *http.Request) {
	raceID := chi.URLParam(r, "id")
	if raceID == "" {
		http.Error(w, "Race ID is required", http.StatusBadRequest)
		return
	}

	race, err := h.raceQueries.FindRaceByID(raceID)
	if err != nil {
		http.Error(w, "Race not found", http.StatusNotFound)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&race)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	err = h.raceQueries.UpdateRace(race)
	if err != nil {
		http.Error(w, "Error updating race", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(race)
}

// CreateRaceHandler creates a new race
// @Summary Create a new race
// @Description Create a new race with the input payload
// @Tags race
// @Accept  json
// @Produce  json
// @Param body body models.Races true "Race object"
// @Success 200 {object} models.Races "Successfully created race"
// @Failure 400 {object} string "Invalid JSON body"
// @Failure 500 {object} string "Error creating race"
// @Router /races [post]
func (h *RaceHandler) CreateRaceHandler(w http.ResponseWriter, r *http.Request) {
	var race models.Races

	err := json.NewDecoder(r.Body).Decode(&race)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	_, err = h.raceQueries.CreateRace(&race)
	if err != nil {
		http.Error(w, "error creating race", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(race)
}

// DeleteRaceHandler deletes a race
// @Summary Delete a race
// @Description Delete a race by ID
// @Tags race
// @Accept  json
// @Produce  json
// @Param id path string true "Race ID"
// @Success 204 "Successfully deleted race"
// @Failure 400 {object} string "Invalid ID supplied"
// @Failure 500 {object} string "Error deleting race"
// @Router /races/{id} [delete]
func (h *RaceHandler) DeleteRaceHandler(w http.ResponseWriter, r *http.Request) {
	raceID := chi.URLParam(r, "id")
	if raceID == "" {
		http.Error(w, "Race ID is required", http.StatusBadRequest)
		return
	}

	err := h.raceQueries.DeleteRace(raceID)
	if err != nil {
		http.Error(w, "Error deleting race", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
