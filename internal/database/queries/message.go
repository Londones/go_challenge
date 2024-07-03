package queries

import (
	"fmt"
	"go-challenge/internal/models"
	"go-challenge/internal/utils"
)

func (s *DatabaseService) SaveMessage(roomID uint, senderID string, content string) (id uint, err error) {
	db := s.s.DB()
	message := models.Message{
		ChatID:   roomID,
		SenderID: senderID,
		Content:  content,
	}
	if err := db.Create(message).Error; err != nil {
		utils.Logger("error", "Message Creation:", "Failed to create message", fmt.Sprintf("Error: %v", err))
		return 0, err
	}
	utils.Logger("info", "Message Creation:", "Message created successfully", fmt.Sprintf("Message ID: %v", message.ID))
	return message.ID, nil
}

func (s *DatabaseService) GetMessagesByRoomID(roomID uint) ([]*models.Message, error) {
	db := s.s.DB()
	var messages []*models.Message
	if err := db.Where("chat_id = ?", roomID).Order("created_at").Find(&messages).Error; err != nil {
		utils.Logger("error", "Get Messages By Room ID:", "Failed to get messages by room ID", fmt.Sprintf("Error: %v", err))
		return nil, err
	}
	utils.Logger("info", "Get Messages By Room ID:", "Messages retrieved successfully", fmt.Sprintf("Messages: %v", messages))
	return messages, nil
}
