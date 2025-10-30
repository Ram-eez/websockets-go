package manager

import (
	"fmt"
	"websockets/models"
)

type Room struct {
	id      string
	clients map[string]*Client
	manager *Manager

	broadcast      chan models.Message
	register       chan *Client
	unregister     chan *Client
	done           chan struct{}
	messageHistory []models.Message
}

func NewRoom(manager *Manager, name string) *Room {
	return &Room{
		id:             name,
		clients:        make(map[string]*Client),
		manager:        manager,
		broadcast:      make(chan models.Message),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		done:           make(chan struct{}),
		messageHistory: make([]models.Message, 0),
	}
}

func (r *Room) Run() {
	for {
		select {

		case client := <-r.register:
			r.clients[client.id] = client
			messages, err := r.manager.repo.GetAllRoomMessages(r.id, 50)
			if err != nil {
				fmt.Println("err could not get room information/messages: ", err)
			}

			for _, msg := range messages {
				client.egress <- *msg
			}
			joinMsg := models.Message{
				Username: "System",
				Message:  client.user.Username + " joined the room ",
				RoomID:   r.id,
			}

			for _, c := range r.clients {
				c.egress <- joinMsg
			}

		case client := <-r.unregister:
			delete(r.clients, client.id)

			leaveMsg := models.Message{
				Username: "System",
				Message:  client.user.Username + " left the room",
				RoomID:   r.id,
			}

			if len(r.clients) == 0 {
				r.manager.RemoveRoom(r.id)
				return
			}
			for _, c := range r.clients {
				c.egress <- leaveMsg
			}

		case msg := <-r.broadcast:
			if msg.Username != "System" {
				if err := r.manager.repo.AddMessage(&msg); err != nil {
					fmt.Println("could not write message:", err)
				}
			}
			for _, client := range r.clients {
				client.egress <- msg
			}
		}
	}
}
