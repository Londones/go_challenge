package queries

import (
	"go-challenge/internal/models"
)

func (s *DatabaseService) CreateChat(chat *models.Chat) (id uint, err error) {
	db := s.s.DB()
	if err := db.Create(chat).Error; err != nil {
		return 0, err
	}
	return chat.ID, nil
}

func (s *DatabaseService) GetOrCreateChat(userID1, userID2 string) (*models.Chat, error) {
	db := s.s.DB()
	chat := &models.Chat{}
	// find the chat where the user is either userID1 or userID2
	if err := db.Where("user1ID = ? AND user2ID = ?", userID1, userID2).First(chat).Error; err != nil {
		// if not found, create a new chat
		chat = &models.Chat{
			User1ID: userID1,
			User2ID: userID2,
		}
		if err := db.Create(chat).Error; err != nil {
			return nil, err
		}
	}
	return chat, nil
}

func (s *DatabaseService) FindChatsByUserID(userid string) ([]*models.Chat, error) {
	db := s.s.DB()
	var chats []*models.Chat
	// find all chats where the user is either userID1 or userID2
	if err := db.Where("userID1 = ? OR userID2 = ?", userid, userid).Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}

// return all chatIDs regardless of the user
func (s *DatabaseService) GetChatIds() ([]uint, error) {
	db := s.s.DB()
	var chatIDs []uint
	if err := db.Table("chats").Pluck("id", &chatIDs).Error; err != nil {
		return nil, err
	}
	return chatIDs, nil
}
