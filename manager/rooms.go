package manager

import (
	"websockets/models"
)

type Room struct {
	id      string
	clients map[string]*Client
	manager *Manager

	broadcast  chan models.Message
	register   chan *Client
	unregister chan *Client
}

func NewRoom(manager *Manager, name string) *Room {
	return &Room{
		id:         name,
		clients:    make(map[string]*Client),
		manager:    manager,
		broadcast:  make(chan models.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (r *Room) Run() {
	for {
		select {

		case client := <-r.register:
			r.clients[client.id] = client

		case client := <-r.unregister:
			delete(r.clients, client.id)

		case msg := <-r.broadcast:
			for _, client := range r.clients {
				client.egress <- msg
			}
		}
	}
}

// func (r *Room) GetRoomHTML() []byte {
// 	tmpl, err := template.ParseFiles("views/newroom.html")
// 	if err != nil {
// 		fmt.Println("templete parsing err: ", err)
// 		return nil
// 	}

// 	var renderedMessage bytes.Buffer

// 	data := map[string]any{
// 		"Room": r.id,
// 	}

// 	if err := tmpl.Execute(&renderedMessage, data); err != nil {
// 		fmt.Println("execution err could not replace : ", err)
// 		return nil
// 	}

// 	fmt.Println("generated HTML with replaced obj: ", renderedMessage.String())

// 	return renderedMessage.Bytes()
// }
