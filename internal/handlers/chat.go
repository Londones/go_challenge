package handlers

import (
	"go-challenge/internal/database/queries"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	queries *queries.DatabaseService
}

func NewChatHandler(queries *queries.DatabaseService) *ChatHandler {
	return &ChatHandler{queries: queries}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	hub         *Hub
	conn        *websocket.Conn
	send        chan []byte
	userID      string
	chatID      uint
	recipientID string
}

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	queries    *ChatHandler
}

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (hub *Hub) Run() {
	for client := range hub.register {
		hub.clients[client] = true
	}
	for client := range hub.unregister {
		if _, ok := hub.clients[client]; ok {
			delete(hub.clients, client)
			close(client.send)
		}
	}
}

func (hub *Hub) sendToRecipient(message []byte, recipientID string, chatID uint) {
	for client := range hub.clients {
		if client.userID == recipientID && client.chatID == chatID {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(hub.clients, client)
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		_, err = c.hub.queries.queries.SaveMessage(c.chatID, c.userID, string(message))
		if err != nil {
			continue
		}

		c.hub.sendToRecipient(message, c.recipientID, c.chatID)
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for message := range c.send {
		c.conn.WriteMessage(websocket.TextMessage, message)
	}
}

func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	userID := r.URL.Query().Get("userID")
	chatIDString := r.URL.Query().Get("chatID")
	chatID, _ := strconv.ParseUint(chatIDString, 10, 32)

	client := &Client{
		hub:         h,
		conn:        conn,
		send:        make(chan []byte),
		userID:      userID,
		chatID:      uint(chatID),
		recipientID: r.URL.Query().Get("recipientID"),
	}
	h.register <- client

	go client.writePump()
	go client.readPump()
}
