package queries

import (
	"errors"
	"fmt"
	"go-challenge/internal/models"
	"strconv"

	"gorm.io/gorm"
)

func (s *DatabaseService) CreateRace(race *models.Races) (uint, error) {
	db := s.s.DB()
	if err := db.Create(race).Error; err != nil {
		return 0, err
	}
	return race.ID, nil
}

func (s *DatabaseService) DeleteRace(id string) error {
	// Vérifiez si l'ID est vide
	if id == "" {
		return fmt.Errorf("l'ID fourni est vide")
	}

	// Convertir l'ID en entier
	raceID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("conversion de l'ID en entier a échoué : %v", err)
	}

	db := s.s.DB()

	var race models.Races
	if err := db.Where("id = ?", raceID).First(&race).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("race avec ID %d introuvable", raceID)
		}
		return err
	}

	if err := db.Delete(&race).Error; err != nil {
		return err
	}

	return nil
}

func (s *DatabaseService) GetAllRace() ([]models.Races, error) {
	db := s.s.DB()

	var races []models.Races
	if err := db.Find(&races).Error; err != nil {
		return nil, err
	}
	return races, nil
}

func (s *DatabaseService) FindRaceByID(id string) (*models.Races, error) {
	// Vérifiez si l'ID est vide
	if id == "" {
		return nil, fmt.Errorf("l'ID fourni est vide")
	}

	var race models.Races

	// Convertir l'ID en entier
	raceID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("conversion de l'ID en entier a échoué : %v", err)
	}

	db := s.s.DB()
	if err := db.Where("ID = ?", raceID).First(&race).Error; err != nil {
		return nil, err
	}
	return &race, nil
}

func (s *DatabaseService) UpdateRace(race *models.Races) error {
	db := s.s.DB()
	if err := db.Save(race).Error; err != nil {
		return err
	}
	return nil
}
