package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"go-challenge/internal/database/queries"
	"go-challenge/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/gorilla/websocket"
)

type RoomHandler struct {
	roomQueries *queries.DatabaseService
	Rooms       *Rooms
	upgrader    websocket.Upgrader
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
		log.Printf("Error loading rooms: %v", err)
		return
	}
	for _, room := range rooms {
		h.Rooms.AddRoom(room)
	}
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
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "error upgrading connection", http.StatusInternalServerError)
		return
	}

	_, claims, _ := jwtauth.FromContext(r.Context())
	userID := claims["userID"].(string)
	roomIDstring := chi.URLParam(r, "roomID")

	roomID, err := strconv.ParseUint(roomIDstring, 10, 64)
	if err != nil {
		http.Error(w, "error parsing room ID", http.StatusBadRequest)
		return
	}

	client := &Client{
		userID: userID,
		conn:   conn,
		send:   make(chan []byte, 256),
	}

	room, err := h.getOrCreateRoom(uint(roomID))
	if err != nil {
		http.Error(w, "room not found", http.StatusNotFound)
		conn.Close()
		return
	}

	room.RegisterClient(client)

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
				log.Printf("error: %v", err)
			}
			break
		}

		_, error := h.roomQueries.SaveMessage(room.roomID, c.userID, string(message))
		if error != nil {
			log.Printf("Error saving message: %v", error)
			break
		}

		room.Broadcast(message)
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for message := range c.send {
		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)

		n := len(c.send)
		for i := 0; i < n; i++ {
			w.Write([]byte{'\n'})
			w.Write(<-c.send)
		}

		if err := w.Close(); err != nil {
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
	userID := claims["userID"].(string)

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
			"user1ID":      room.User1ID,
			"user2ID":      room.User2ID,
			"annonceID":    room.AnnonceID,
			"annonceTitle": annonce.Title,
		}

		modifiedRooms = append(modifiedRooms, modifiedRoom)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(modifiedRooms)

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
		http.Error(w, "room ID is required", http.StatusBadRequest)
		return
	}

	roomIDToInt, err := strconv.ParseUint(roomID, 10, 64)
	if err != nil {
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
}
