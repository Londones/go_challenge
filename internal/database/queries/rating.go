package queries

import (
	"errors"
	"fmt"
	"go-challenge/internal/models"

	"gorm.io/gorm"
)

func (s *DatabaseService) CreateRating(rating *models.Rating) (id uint, err error) {
	db := s.s.DB()
	if err := db.Create(rating).Error; err != nil {
		return 0, err
	}
	return rating.ID, nil
}

func (s *DatabaseService) FindRatingByID(id string) (*models.Rating, error) {
	db := s.s.DB()
	var rating models.Rating
	if err := db.Where("ID = ?", id).First(&rating).Error; err != nil {
		return nil, err
	}
	return &rating, nil
}

func (s *DatabaseService) UpdateRatingComment(id string, comment string) (*models.Rating, error) {
	db := s.s.DB()

	var rating models.Rating
	if err := db.Where("id = ?", id).First(&rating).Error; err != nil {
		return nil, err
	}

	rating.Comment = comment

	if err := db.Save(&rating).Error; err != nil {
		return nil, err
	}

	return &rating, nil
}

func (s *DatabaseService) DeleteRating(id string) error {
	db := s.s.DB()

	var rating models.Rating
	if err := db.Where("id = ?", id).First(&rating).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("rating with ID %s not found", id)
		}
		return err
	}

	if err := db.Delete(&rating).Error; err != nil {
		return err
	}

	return nil
}

func (s *DatabaseService) GetAllRatings() ([]models.Rating, error) {
	db := s.s.DB()
	var ratings []models.Rating
	if err := db.Find(&ratings).Error; err != nil {
		return nil, err
	}
	return ratings, nil
}

func (s *DatabaseService) GetUserRatings(userID uint) ([]models.Rating, error) {
	db := s.s.DB()
	var ratings []models.Rating
	if err := db.Where("user_id = ?", userID).Find(&ratings).Error; err != nil {
		return nil, err
	}
	return ratings, nil
}

func (s *DatabaseService) GetAnnonceRatings(annonceID string) ([]models.Rating, error) {
	db := s.s.DB()
	var ratings []models.Rating
	if err := db.Where("annonce_id = ?", annonceID).Find(&ratings).Error; err != nil {
		return nil, err
	}
	return ratings, nil
}

func (s *DatabaseService) UpdateRating(id string, mark int8, comment string) (*models.Rating, error) {
	db := s.s.DB()

	var rating models.Rating
	if err := db.Where("id = ?", id).First(&rating).Error; err != nil {
		return nil, err
	}

	rating.Mark = mark
	if comment != "" {
		rating.Comment = comment
	}

	if err := db.Save(&rating).Error; err != nil {
		return nil, err
	}

	return &rating, nil
}
