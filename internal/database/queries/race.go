package queries

import (
	"errors"
	"fmt"
	"go-challenge/internal/models"
	"strconv"

	"gorm.io/gorm"
)

func (s *DatabaseService) CreateRace(race *models.Races) (id uint, err error) {
	fmt.Println("creation d'une race de chat")
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

func (s *DatabaseService) FindRaceByID(id string) (race models.Races, err error) {
	// Vérifiez si l'ID est vide
	if id == "" {
		return models.Races{}, fmt.Errorf("l'ID fourni est vide")
	}

	// Convertir l'ID en entier
	raceID, err := strconv.Atoi(id)
	if err != nil {
		return models.Races{}, fmt.Errorf("conversion de l'ID en entier a échoué : %v", err)
	}

	db := s.s.DB()
	if err := db.Where("ID = ?", raceID).First(&race).Error; err != nil {
		return models.Races{}, err
	}
	return race, nil
}

func (s *DatabaseService) DeleteRaceById(id string) error {
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
			return fmt.Errorf("Race avec ID %d introuvable", raceID)
		}
		return err
	}

	if err := db.Delete(&race).Error; err != nil {
		return err
	}

	return nil
}

func (s *DatabaseService) UpdateRace(race models.Races) error {
	db := s.s.DB()

	var existingRace models.Races
	if err := db.Where("id = ?", race.ID).First(&existingRace).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("Race avec ID %d introuvable", race.ID)
		}
		return err
	}

	if err := db.Model(&existingRace).Updates(race).Error; err != nil {
		return err
	}

	return nil
}
