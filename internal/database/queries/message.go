package queries

import (
	"go-challenge/internal/models"
)

func (s *DatabaseService) SaveMessage(roomID uint, senderID string, content string) (id uint, err error) {
	db := s.s.DB()
	message := models.Message{
		ChatID:   roomID,
		SenderID: senderID,
		Content:  content,
	}
	if err := db.Create(message).Error; err != nil {
		return 0, err
	}
	return message.ID, nil
}

func (s *DatabaseService) GetMessagesByRoomID(roomID uint) ([]*models.Message, error) {
	db := s.s.DB()
	var messages []*models.Message
	if err := db.Where("chat_id = ?", roomID).Order("created_at").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
