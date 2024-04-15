package sockethandler

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
)

type Room struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Clients   map[string]*Client
	Broadcast chan Message
}

type RoomHelper interface {
	Get() *Room
	GetId() string
	GetName() string
	Write()
	Read(client *Client)
	RemoveClient(client *Client)
	EmitAll(msg Message)
	Emit(msg Message)
}

func NewRoom(id string, name string) RoomHelper {
	return &Room{
		Id:        id,
		Name:      name,
		Clients:   make(map[string]*Client),
		Broadcast: make(chan Message),
	}
}

func (r *Room) Get() *Room {
	return r
}

func (r *Room) GetId() string {
	return r.Id
}

func (r *Room) GetName() string {
	return r.Name
}

func (r *Room) Write() {
	for {
		select {
		case msg := <-r.Broadcast:
			client := r.Clients[msg.ClientId]
			if client != nil && client.Conn != nil {
				if err := client.Write(msg); !errors.Is(err, nil) {
					log.Printf("error occurred: %v", err)
				}
			}
		}
	}
}

func (r *Room) Read(client *Client) {
	for {
		var req Message
		err := client.Read(&req)

		if err != nil {
			log.Printf("Invalid data sent from client: %v", err.Error())
			r.RemoveClient(client)
			break
		} else {
			r.EmitAll(Message{
				Event:      req.Event,
				Data:       req.Data,
				ClientId:   client.Id,
				ClientName: client.Name,
				Sender: Sender{
					Id:   client.Id,
					Name: client.Name,
				},
			})
		}
	}
}

func (r *Room) Emit(msg Message) {
	r.Broadcast <- msg
}

func (r *Room) EmitAll(msg Message) {
	for _, client := range r.Clients {
		msg.ClientId = client.Id
		msg.ClientName = client.Name
		r.Emit(msg)
	}
}

func (r *Room) RemoveClient(client *Client) {

	r.EmitAll(Message{
		Event:      "chat:message::event",
		Data:       fmt.Sprintf("%s left", client.Name),
		ClientId:   client.Id,
		ClientName: client.Name,
		Sender: Sender{
			Id:   client.Id,
			Name: client.Name,
		},
	})

	delete(r.Clients, client.Id)
}
