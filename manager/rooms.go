package manager

import (
	"websockets/models"
)

type Room struct {
	id      string
	clients map[string]*Client
	manager *Manager

	broadcast      chan models.Message
	register       chan *Client
	unregister     chan *Client
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
		messageHistory: make([]models.Message, 0),
	}
}

func (r *Room) Run() {
	for {
		select {

		case client := <-r.register:
			r.clients[client.id] = client
			for _, msg := range r.messageHistory {
				client.egress <- msg
			}

		case client := <-r.unregister:
			delete(r.clients, client.id)

		case msg := <-r.broadcast:
			r.messageHistory = append(r.messageHistory, msg)
			for _, client := range r.clients {
				client.egress <- msg
			}
		}
	}
}
