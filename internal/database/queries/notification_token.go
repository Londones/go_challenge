package queries

import (
	"go-challenge/internal/models"
)

func (s *DatabaseService) CreateNotificationToken(notificationToken *models.NotificationToken) error {
    db := s.s.DB()
    if err := db.Where(models.NotificationToken{UserID: notificationToken.UserID}).Assign(models.NotificationToken{Token: notificationToken.Token}).FirstOrCreate(notificationToken).Error; err != nil {
        return err
    }
    return nil
}

func (s *DatabaseService) DeleteNotificationTokenForUser(userID string) error {
	db := s.s.DB()
	if err := db.Where("user_id = ?", userID).Delete(&models.NotificationToken{}).Error; err != nil {
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

