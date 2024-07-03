package queries

import (
	"go-challenge/internal/models"
)

func (s *DatabaseService) CreateRoom(room *models.Room) (id uint, err error) {
	db := s.s.DB()
	if err := db.Create(room).Error; err != nil {
		return 0, err
	}
	return room.ID, nil
}

func (s *DatabaseService) GetOrCreateRoom(userID1, userID2 string) (*models.Room, error) {
	db := s.s.DB()
	room := &models.Room{}
	// find the room where the user is either userID1 or userID2
	if err := db.Where("user1ID = ? AND user2ID = ?", userID1, userID2).First(room).Error; err != nil {
		// if not found, create a new room
		room = &models.Room{
			User1ID: userID1,
			User2ID: userID2,
		}
		if err := db.Create(room).Error; err != nil {
			return nil, err
		}
	}
	return room, nil
}

func (s *DatabaseService) FindRoomsByUserID(userid string) ([]*models.Room, error) {
	db := s.s.DB()
	var rooms []*models.Room
	// find all rooms where the user is either userID1 or userID2
	if err := db.Where("userID1 = ? OR userID2 = ?", userid, userid).Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

// return all roomIDs regardless of the user
func (s *DatabaseService) GetRoomIds() ([]uint, error) {
	db := s.s.DB()
	var roomIDs []uint
	if err := db.Table("rooms").Pluck("id", &roomIDs).Error; err != nil {
		return nil, err
	}
	return roomIDs, nil
}

func (s *DatabaseService) GetRoomByID(roomID uint) (*models.Room, error) {
	db := s.s.DB()
	room := &models.Room{}
	if err := db.First(room, roomID).Error; err != nil {
		return nil, err
	}
	return room, nil
}

func (s *DatabaseService) GetRooms() ([]*models.Room, error) {
	db := s.s.DB()
	var rooms []*models.Room
	if err := db.Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}
