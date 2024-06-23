package queries

import (
	"go-challenge/internal/models"
)

func (s *DatabaseService) SaveMessage(chatID uint, senderID string, content string) (id uint, err error) {
	db := s.s.DB()
	message := models.Message{
		ChatID:   chatID,
		SenderID: senderID,
		Content:  content,
	}
	if err := db.Create(message).Error; err != nil {
		return 0, err
	}
	return message.ID, nil
}
