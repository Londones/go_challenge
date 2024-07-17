package queries

import (
	"fmt"
	"go-challenge/internal/models"
	"go-challenge/internal/utils"
	"gorm.io/gorm"
)

func (s *DatabaseService) SaveMessage(roomID uint, senderID string, content string) (*models.Message, string, error) {
	db := s.s.DB()
	message := models.Message{
		RoomID:   roomID,
		SenderID: senderID,
		Content:  content,
	}
	if err := db.Create(&message).Error; err != nil {
		utils.Logger("error", "Message Creation:", "Failed to create message", fmt.Sprintf("Error: %v", err))
		return nil, "", err
	}
	utils.Logger("info", "Message Creation:", "Message created successfully", fmt.Sprintf("Message ID: %v", message.ID))
	

	user, err := s.FindUserByID(senderID)
	if err != nil {
		utils.Logger("error", "Message Creation:", "Failed to get user by ID", fmt.Sprintf("Error: %v", err))
		return nil, "", err
	}
	return &message, user.Name, nil


}

func (s *DatabaseService) GetMessagesByRoomID(roomID uint) ([]*models.Message, error) {
	db := s.s.DB()
	var messages []*models.Message
	if err := db.Where("room_id = ?", roomID).Order("created_at").Find(&messages).Error; err != nil {
		utils.Logger("error", "Get Messages By Room ID:", "Failed to get messages by room ID", fmt.Sprintf("Error: %v", err))
		return nil, err
	}
	utils.Logger("info", "Get Messages By Room ID:", "Messages retrieved successfully", fmt.Sprintf("Messages: %v", messages))
	return messages, nil
}

func (s *DatabaseService) MarkMessagesAsRead(roomID uint, userID string) error {
	db := s.s.DB()
	if err := db.Model(&models.Message{}).Where("room_id = ? AND sender_id != ? AND is_read = false", roomID, userID).Update("is_read", true).Error; err != nil {
		utils.Logger("error", "Mark Messages As Read:", "Failed to mark messages as read", fmt.Sprintf("Error: %v", err))
		return err
	}
	utils.Logger("info", "Mark Messages As Read:", "Messages marked as read successfully", fmt.Sprintf("Room ID: %v, Sender ID: %v", roomID, userID))
	return nil
}

func (s *DatabaseService) GetLatestMessageByRoomID(roomID uint) (*LatestMessageResponse, error) {
	db := s.s.DB()
	var message models.Message
	result := db.Where("room_id = ?", roomID).
		Order("created_at DESC").
		First(&message)
	if result.Error != nil {
		if result.Error.Error() == gorm.ErrRecordNotFound.Error() {
			return &LatestMessageResponse{Message: nil}, nil
		}
		return nil, result.Error
	}
	return &LatestMessageResponse{Message: &message}, nil
}

type LatestMessageResponse struct {
	Message *models.Message `json:"message"`
}
