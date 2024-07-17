package queries

import (
	"fmt"
	"go-challenge/internal/models"
	"go-challenge/internal/utils"
)

func (s *DatabaseService) CreateRoom(room *models.Room) (id uint, err error) {
	db := s.s.DB()
	if err := db.Create(room).Error; err != nil {
		utils.Logger("error", "Room Creation:", "Failed to create room", fmt.Sprintf("Error: %v", err))
		return 0, err
	}
	utils.Logger("info", "Room Creation:", "Room created successfully", fmt.Sprintf("Room ID: %v", room.ID))
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
			utils.Logger("error", "Room Get or Creation:", "Failed to get or create room", fmt.Sprintf("Error: %v", err))
			return nil, err
		}
	}
	utils.Logger("info", "Room Get or Creation:", "Room retrieved or created successfully", fmt.Sprintf("Room ID: %v", room.ID))
	return room, nil
}

func (s *DatabaseService) FindRoomsByUserID(userid string) ([]*models.Room, error) {
	db := s.s.DB()
	var rooms []*models.Room
	// find all rooms where the user is either userID1 or userID2
	if err := db.Where("user1_ID = ? OR user2_ID = ?", userid, userid).Find(&rooms).Error; err != nil {
		utils.Logger("error", "Find Rooms By User ID:", "Failed to find rooms by user ID", fmt.Sprintf("Error: %v", err))
		return nil, err
	}
	utils.Logger("info", "Find Rooms By User ID:", "Rooms found successfully", fmt.Sprintf("Rooms: %v", rooms))
	return rooms, nil
}

// return all roomIDs regardless of the user
func (s *DatabaseService) GetRoomIds() ([]uint, error) {
	db := s.s.DB()
	var roomIDs []uint
	if err := db.Table("rooms").Pluck("id", &roomIDs).Error; err != nil {
		utils.Logger("error", "Get Room IDs:", "Failed to get room IDs", fmt.Sprintf("Error: %v", err))
		return nil, err
	}
	utils.Logger("info", "Get Room IDs:", "Room IDs retrieved successfully", fmt.Sprintf("Room IDs: %v", roomIDs))
	return roomIDs, nil
}

func (s *DatabaseService) GetRoomByID(roomID uint) (*models.Room, error) {
	db := s.s.DB()
	room := &models.Room{}
	if err := db.First(room, roomID).Error; err != nil {
		utils.Logger("error", "Get Room By ID:", "Failed to get room by ID", fmt.Sprintf("Error: %v", err))
		return nil, err
	}
	utils.Logger("info", "Get Room By ID:", "Room retrieved successfully", fmt.Sprintf("Room: %v", room))
	return room, nil
}

func (s *DatabaseService) GetRooms() ([]*models.Room, error) {
	db := s.s.DB()
	var rooms []*models.Room
	if err := db.Find(&rooms).Error; err != nil {
		utils.Logger("error", "Get Rooms:", "Failed to get rooms", fmt.Sprintf("Error: %v", err))
		return nil, err
	}
	utils.Logger("info", "Get Rooms:", "Rooms retrieved successfully", fmt.Sprintf("Rooms: %v", rooms))
	return rooms, nil
}
