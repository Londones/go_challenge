package queries

import (
	"go-challenge/internal/models"
)

func (s *DatabaseService) CreateNotificationToken(notificationToken *models.NotificationToken) error {
	db := s.s.DB()
	if err := db.Create(notificationToken).Error; err != nil {
		return err
	}
	return nil
}

func (s *DatabaseService) DeleteNotificationToken(id int) error {
	db := s.s.DB()
	if err := db.Where("id = ?", id).Delete(&models.NotificationToken{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *DatabaseService) GetNotificationTokenByUserID(userID string) (*models.NotificationToken, error) {
	db := s.s.DB()
	var notificationToken models.NotificationToken
	if err := db.Where("user_id = ?", userID).First(&notificationToken).Error; err != nil {
		return nil, err
	}
	return &notificationToken, nil
}

