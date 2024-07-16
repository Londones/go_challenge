package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"go-challenge/internal/database/queries"
	"go-challenge/internal/models"
	"go-challenge/internal/config"
	"go-challenge/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/gorilla/websocket"
)

type RoomHandler struct {
	roomQueries *queries.DatabaseService
	Rooms       *Rooms
	upgrader    websocket.Upgrader
}

type MessageJSON struct {
	Content   string    `json:"Content"`
	SenderID  string    `json:"SenderID"`
	CreatedAt time.Time `json:"CreatedAt"`
}

func NewRoomHandler(roomQueries *queries.DatabaseService) *RoomHandler {
	return &RoomHandler{
		roomQueries: roomQueries,
		Rooms:       NewRooms(),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

type Client struct {
	userID string
	conn   *websocket.Conn
	send   chan []byte
}

type Room struct {
	roomID  uint
	clients map[string]*Client
	mu      sync.Mutex
}

type Rooms struct {
	rooms map[uint]*Room
	mu    sync.RWMutex
}

func NewRooms() *Rooms {
	return &Rooms{
		rooms: make(map[uint]*Room),
	}
}

func (h *RoomHandler) LoadRooms() {
	rooms, err := h.roomQueries.GetRooms()
	if err != nil {
		utils.Logger("error", "Load Rooms:", "Failed to load rooms", fmt.Sprintf("Error: %v", err))
		log.Printf("Error loading rooms: %v", err)
		return
	}
	for _, room := range rooms {
		h.Rooms.AddRoom(room)
	}
	utils.Logger("info", "Load Rooms:", "Rooms loaded successfully", "")
}

func (r *Rooms) AddRoom(room *models.Room) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rooms[room.ID] = &Room{
		roomID:  room.ID,
		clients: make(map[string]*Client),
	}
}

// HandleWebSocket godoc
// @Summary Connect to WebSocket for real-time chat
// @Description Establishes a WebSocket connection to a specified chat room for real-time communication.
// @Param roomID path int true "The ID of the chat room to connect to"
// @Param Authorization header string true "Bearer token for authentication"
// @Success 101 "Switching Protocols"
// @Failure 400 "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Security BearerAuth
// @Router /ws/{roomID} [get]
func (h *RoomHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {

	_, claims, _ := jwtauth.FromContext(r.Context())
	userID := claims["id"].(string)
	roomIDstring := chi.URLParam(r, "roomID")

	roomID, err := strconv.ParseUint(roomIDstring, 10, 64)
	if err != nil {
		utils.Logger("error", "Handle WebSocket:", "Failed to parse room ID", fmt.Sprintf("Error: %v", err))
		http.Error(w, "error parsing room ID", http.StatusBadRequest)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.Logger("error", "Handle WebSocket:", "Failed to upgrade connection", fmt.Sprintf("Error: %v", err))
		http.Error(w, "error upgrading connection", http.StatusInternalServerError)
		return
	}

	client := &Client{
		userID: userID,
		conn:   conn,
		send:   make(chan []byte, 256),
	}

	room, err := h.getOrCreateRoom(uint(roomID))
	if err != nil {
		utils.Logger("error", "Handle WebSocket:", "Failed to get or create room", fmt.Sprintf("Error: %v", err))
		http.Error(w, "room not found", http.StatusNotFound)
		conn.Close()
		return
	}

	room.RegisterClient(client)

	utils.Logger("info", "Handle WebSocket:", "Client connected to room", fmt.Sprintf("Room ID: %v, User ID: %v", roomID, userID))

	go client.writePump()
	go client.readPump(room, h)
}

func (r *Rooms) GetRoom(roomID uint) *Room {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.rooms[roomID]
}

func (h *RoomHandler) getOrCreateRoom(roomID uint) (*Room, error) {
	h.Rooms.mu.Lock()
	defer h.Rooms.mu.Unlock()

	if room, exists := h.Rooms.rooms[roomID]; exists {
		return room, nil
	}

	// Room doesn't exist, fetch from database
	dbRoom, err := h.roomQueries.GetRoomByID(roomID)
	if err != nil {
		return nil, err
	}

	room := &Room{
		roomID:  dbRoom.ID,
		clients: make(map[string]*Client),
	}
	h.Rooms.rooms[roomID] = room

	return room, nil
}

func (r *Room) RegisterClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[client.userID] = client
}

func (r *Room) UnregisterClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.clients[client.userID]; ok {
		delete(r.clients, client.userID)
		close(client.send)
	}
}

func (r *Room) Broadcast(message []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, client := range r.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(r.clients, client.userID)
		}
	}
}

func (c *Client) readPump(room *Room, h *RoomHandler) {
	defer func() {
		room.UnregisterClient(c)
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				utils.Logger("error", "Read Pump:", "Unexpected close error", fmt.Sprintf("Error: %v", err))
				fmt.Println("error: ", err)
			}
			break
		}

		createdMessage, error := h.roomQueries.SaveMessage(room.roomID, c.userID, string(message))
		if error != nil {
			log.Printf("Error saving message: %v", error)
			break
		}

		for k, _ := range room.clients {
			if k != c.userID {
				notificationToken, error := h.roomQueries.GetNotificationTokenByUserID(k)
				if error != nil {
					fmt.Println("Error getting notification token: ", error)
				} else if notificationToken.Token != "" {
					fmt.Println("Sending notification....", notificationToken.Token)
					payload := make(map[string]string)
					payload["RoomID"] = strconv.FormatUint(uint64(room.roomID), 10)
					SendToToken(config.GetFirebaseApp(), notificationToken.Token, createdMessage.Content, createdMessage.SenderID, payload)
				}
			}
		}
		
		// disconnectedUserID, err := h.GetDisconnectedUser(room.roomID)
		// if err != nil {
		// 	log.Printf("Error getting disconnected user: %v", err)
		// } else {
		// 	notificationToken, error := h.roomQueries.GetNotificationTokenByUserID(disconnectedUserID)
		// 	if error != nil {
		// 		log.Printf("Error getting notification token: %v", error)
		// 	} else if notificationToken.Token != "" {
		// 		log.Printf("Sending notification....%v", notificationToken.Token)
		// 		SendToToken(config.GetFirebaseApp(), notificationToken.Token, createdMessage.Content, createdMessage.SenderID)
		// 	}
		// }
		

		/*jsonMessage := MessageJSON{
			Content:   createdMessage.Content,
			SenderID:  createdMessage.SenderID,
			CreatedAt: createdMessage.CreatedAt.Format(time.RFC3339),
		}*/

		message, err = json.Marshal(createdMessage)
		if err != nil {
			utils.Logger("error", "Read Pump:", "Failed to marshal message", fmt.Sprintf("Error: %v", err))
			log.Printf("error: %v", err)
			continue
		}

		room.Broadcast(message)
		utils.Logger("info", "Read Pump:", "Message broadcasted", fmt.Sprintf("Room ID: %v, User ID: %v", room.roomID, c.userID))
	}
}


func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for message := range c.send {
		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			utils.Logger("error", "Write Pump:", "Failed to get writer", fmt.Sprintf("Error: %v", err))
			return
		}
		w.Write(message)
		utils.Logger("info", "Write Pump:", "Message sent", fmt.Sprintf("Message: %v", message))
		utils.Logger("info", "Write Pump:", "Message sent to", fmt.Sprintf("User ID: %v %v", c.userID, message))

		n := len(c.send)
		for i := 0; i < n; i++ {
			w.Write([]byte{'\n'})
			w.Write(<-c.send)
		}

		utils.Logger("info", "Write Pump:", "Messages sent", fmt.Sprintf("Number of messages: %v", n))

		if err := w.Close(); err != nil {
			utils.Logger("error", "Write Pump:", "Failed to close writer", fmt.Sprintf("Error: %v", err))
			return
		}
	}
}

// GetUserRooms godoc
// @Summary Get all rooms for a user
// @Description Get all rooms for a user
// @Tags rooms
// @Accept json
// @Produce json
// @Success 200 {array} models.Room "rooms for the user"
// @Failure 500 {string} string "error getting rooms"
// @Router /rooms [get]
func (h *RoomHandler) GetUserRooms(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	userID := claims["id"].(string)

	rooms, err := h.roomQueries.FindRoomsByUserID(userID)
	if err != nil {
		http.Error(w, "error getting rooms", http.StatusInternalServerError)
		return
	}

	var modifiedRooms []map[string]interface{}

	// add the annonce title to each room object
	for _, room := range rooms {
		annonce, err := h.roomQueries.FindAnnonceByID(room.AnnonceID)
		if err != nil {
			http.Error(w, "error getting annonce", http.StatusInternalServerError)
			return
		}

		modifiedRoom := map[string]interface{}{
			"id":           room.ID,
			"user1ID":      string(room.User1ID),
			"user2ID":      string(room.User2ID),
			"annonceID":    annonce.ID,
			"annonceTitle": string(annonce.Title),
		}

		modifiedRooms = append(modifiedRooms, modifiedRoom)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(modifiedRooms); err != nil {
		http.Error(w, "error encoding rooms to JSON", http.StatusInternalServerError)
	}
}

// GetRoomMessages godoc
// @Summary Get all messages for a room in order
// @Description Get all messages for a room in order
// @Tags rooms
// @Accept json
// @Produce json
// @Param roomID path string true "ID of the room"
// @Success 200 {array} models.Message "messages for the room"
// @Failure 400 {string} string "room ID is required"
// @Failure 404 {string} string "room not found"
// @Failure 500 {string} string "error getting messages"
// @Router /rooms/{roomID} [get]
func (h *RoomHandler) GetRoomMessages(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")
	if roomID == "" {
		utils.Logger("error", "Get Room Messages:", "Room ID is required", "")
		http.Error(w, "room ID is required", http.StatusBadRequest)
		return
	}

	roomIDToInt, err := strconv.ParseUint(roomID, 10, 64)
	if err != nil {
		utils.Logger("error", "Get Room Messages:", "Failed to parse room ID", fmt.Sprintf("Error: %v", err))
		http.Error(w, "error parsing room ID", http.StatusBadRequest)
		return
	}

	messages, err := h.roomQueries.GetMessagesByRoomID(uint(roomIDToInt))
	if err != nil {
		http.Error(w, "error getting messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
	utils.Logger("info", "Get Room Messages:", "Messages retrieved successfully", fmt.Sprintf("Room ID: %v", roomID))
}


func (h *RoomHandler) GetDisconnectedUser(roomID uint) (string, error) {
    room := h.Rooms.GetRoom(roomID)
    if room == nil {
        return "", fmt.Errorf("room not found")
    }

    room.mu.Lock()
    defer room.mu.Unlock()

    dbRoom, err := h.roomQueries.GetRoomByID(roomID)
    if err != nil {
        return "", err
    }

    // Check user not connected
    if _, user1Connected := room.clients[dbRoom.User1ID]; !user1Connected {
        return dbRoom.User1ID, nil
    }
    if _, user2Connected := room.clients[dbRoom.User2ID]; !user2Connected {
        return dbRoom.User2ID, nil
    }

    return "", fmt.Errorf("both users are connected")
}

