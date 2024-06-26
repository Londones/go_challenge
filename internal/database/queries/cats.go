package queries

import (
	"errors"
	"fmt"
	"go-challenge/internal/models"
	"time"

	"gorm.io/gorm"
)

func (s *DatabaseService) CreateCat(cat *models.Cats) (id uint, err error) {

	fmt.Println("creation du chat")
	db := s.s.DB()
	if err := db.Create(cat).Error; err != nil {
		return 0, err
	}
	return cat.ID, nil
}

func (s *DatabaseService) FindCatByID(id string) (*models.Cats, error) {

	db := s.s.DB()
	var cat models.Cats
	if err := db.Where("ID = ?", id).First(&cat).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

func (s *DatabaseService) DeleteCat(id string) error {
	db := s.s.DB()

	var cat models.Cats
	if err := db.Where("id = ?", id).First(&cat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("cat with ID %s not found", id)
		}
		return err
	}

	if err := db.Delete(&cat).Error; err != nil {
		return err
	}

	return nil
}

func (s *DatabaseService) GetAllCats() ([]models.Cats, error) {
	db := s.s.DB()
	var cats []models.Cats
	if err := db.Find(&cats).Error; err != nil {
		return nil, err
	}
	return cats, nil
}

func (s *DatabaseService) DeleteCatByID(id string) error {
	db := s.s.DB()

	var cat models.Cats
	if err := db.Where("id = ?", id).First(&cat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("cat with ID %s not found", id)
		}
		return err
	}

	if err := db.Delete(&cat).Error; err != nil {
		return err
	}

	return nil
}

func (s *DatabaseService) UpdateCat(cat *models.Cats) error {
	db := s.s.DB()

	var existingCat models.Cats
	if err := db.Where("id = ?", cat.ID).First(&existingCat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("cat with ID %d not found", cat.ID)
		}
		return err
	}

	if err := db.Model(&existingCat).Updates(cat).Error; err != nil {
		return err
	}

	return nil
}

func (s *DatabaseService) GetCatByFilters(raceId int, age int, sex bool) ([]models.Cats, error) {
	var cats []models.Cats
	var birthDate time.Time
	db := s.s.DB()

	birthDate = time.Now().AddDate(-age, 0, 0)

	if err := db.Where("Sexe = ?", &sex).Where("RaceID = ?", &raceId).Where("BirthDate <= ?", &birthDate).Find(&cats).Error; err != nil {
		return nil, err
	}

	return cats, nil
}
