package queries

import (
	"go-challenge/internal/models"
)

func (s *DatabaseService) CreateRating(rating *models.Rating) (uint, error) {
	db := s.s.DB()
	if err := db.Create(rating).Error; err != nil {
		return 0, err
	}
	return rating.ID, nil
}

func (s *DatabaseService) UpdateRating(rating *models.Rating) error {
	db := s.s.DB()
	return db.Save(rating).Error
}

func (s *DatabaseService) DeleteRating(id string) error {
	db := s.s.DB()
	return db.Delete(&models.Rating{}, "id = ?", id).Error
}

func (s *DatabaseService) FindRatingByID(id string) (*models.Rating, error) {
	db := s.s.DB()
	var rating models.Rating
	if err := db.Where("id = ?", id).First(&rating).Error; err != nil {
		return nil, err
	}
	return &rating, nil
}

func (s *DatabaseService) GetAllRatings() ([]models.Rating, error) {
	db := s.s.DB()
	var ratings []models.Rating
	if err := db.Find(&ratings).Error; err != nil {
		return nil, err
	}
	return ratings, nil
}

func (s *DatabaseService) GetUserRatings(userID string) ([]models.Rating, error) {
	db := s.s.DB()
	var ratings []models.Rating
	if err := db.Where("user_id = ?", userID).Find(&ratings).Error; err != nil {
		return nil, err
	}
	return ratings, nil
}

func (s *DatabaseService) GetAuthorRatings(authorID string) ([]models.Rating, error) {
	db := s.s.DB()
	var ratings []models.Rating
	if err := db.Where("author_id = ?", authorID).Find(&ratings).Error; err != nil {
		return nil, err
	}
	return ratings, nil
}
