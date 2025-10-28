package manager

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"websockets/models"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	id         string
	connection *websocket.Conn
	user       *models.User
	manager    *Manager
	// egress channel is an unbuffered channel which is used to avoid concurrent writes on the websocket conn
	egress    chan models.Message
	closeOnce sync.Once
}

func NewClient(conn *websocket.Conn, manager *Manager, user *models.User) *Client {
	NewUUID := uuid.New()
	return &Client{
		id:         NewUUID.String(),
		connection: conn,
		user:       user,
		manager:    manager,
		egress:     make(chan models.Message),
	}
}

func (c *Client) readMessages() {
	defer func() {
		c.manager.UnregisterEverywhere(c)
	}()
	for {
		messageType, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("WebSocket error:", err)
			}
			break
		}

		// Log raw message for debugging
		fmt.Printf("Raw message received: %s\n", string(payload))

		var msg models.Message
		if err := json.Unmarshal(payload, &msg); err != nil {
			fmt.Println("JSON unmarshalling error:", err)
			fmt.Println("Raw payload was:", string(payload))
			continue
		}

		// Set username from authenticated user
		msg.Username = c.user.Username

		// Determine room ID
		roomID := msg.RoomID
		if roomID == "" {
			roomID = "lobby" // Default to lobby instead of user's personal lobby
		}

		fmt.Printf("Message from %s in room %s: %s\n", msg.Username, roomID, msg.Message)

		r := c.manager.GetOrCreateRoom(roomID)
		r.broadcast <- msg

		fmt.Printf("Message type: %d\n", messageType)
	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.UnregisterEverywhere(c)
	}()

	for msg := range c.egress {
		// Log the message being sent for debugging
		htmlContent := msg.GetMessageHTML()
		fmt.Printf("Sending HTML to client: %s\n", string(htmlContent))

		if err := c.manager.repo.AddMessage(&msg); err != nil {
			fmt.Println("could not write the message into repo: ", err)
		}

		if err := c.connection.WriteMessage(
			websocket.TextMessage,
			htmlContent,
		); err != nil {
			fmt.Printf("Failed to send message to client %s: %v\n", c.id, err)
			return
		}
	}
}
