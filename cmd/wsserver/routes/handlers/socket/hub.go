package sockethandler

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type SenderOptions struct {
	AccountId string
	UserId    string
	Type      string
}

type ReadRequest struct {
	Event   string                 `json:"event"`
	Message map[string]interface{} `json:"message"`
}

type WriteResponse struct {
	Event   string
	Message map[string]interface{}
	Sender  SenderOptions
}
type Hub struct {
	Id          string
	ClientId    string
	Clients     map[string]*Client
	Broadcast   chan WriteResponse
	SendToWrite chan WriteResponse
}

type HubHelper interface {
	Get() *Hub
	Write()
	Read(client *Client)
	RemoveClient(client *Client, event string)
	TestSubscribeLogs(client *Client, event string)
}

func NewConnHub(Id string, clientId string) HubHelper {
	return &Hub{
		Id:          Id,
		ClientId:    clientId,
		Clients:     make(map[string]*Client),
		Broadcast:   make(chan WriteResponse),
		SendToWrite: make(chan WriteResponse),
	}
}

func (h *Hub) Get() *Hub {
	return h
}

func (h *Hub) Write() {
	for {
		select {
		case message := <-h.Broadcast:
			for _, client := range h.Clients {
				if err := client.WsConn.WriteJSON(message.Message); !errors.Is(err, nil) {
					log.Printf("error occurred: %v", err)
				}
			}
		case privateMessage := <-h.SendToWrite:
			client := h.Clients[privateMessage.Sender.UserId]
			if client != nil && client.WsConn != nil {
				if err := client.WsConn.WriteJSON(privateMessage.Message); !errors.Is(err, nil) {
					log.Printf("error occurred: %v", err)
				}
			}
		}
	}
}

func (h *Hub) Read(client *Client) {
	for {
		var req map[string]interface{}
		err := client.WsConn.ReadJSON(&req)

		if err != nil {
			log.Printf("Invalid data sent from client: %v", err.Error())
			h.RemoveClient(client, "")
			break
		} else {
			if event, ok := req["event"].(string); ok {
				if event == "ping" {
					h.SendToWrite <- WriteResponse{
						Message: map[string]interface{}{
							"ping": "pong",
						},
						Sender: SenderOptions{
							AccountId: h.Id,
							UserId:    client.Id,
						},
					}
				} else if event == "unsubscribe" {
					client.RemoveTimer(event)
				} else if event == "subscribe" {
					client.RemoveTimer(event)
					// go client.DataConsumer.StartConsumer(req, client.DataProducer)
				}
			}
		}
	}
}

func (h *Hub) RemoveClient(client *Client, event string) {
	if event != "" {
		client.RemoveTimer(event)
	}
	//client.DataProducer.CloseProducer()
	//client.DataConsumer.CloseConsumer()
	delete(h.Clients, client.Id)
}

func (h *Hub) Dispatch(req WriteResponse) {
	h.SendToWrite <- req
}

func (h *Hub) TestSubscribeLogs(client *Client, event string) {
	ticker := time.NewTicker(3 * time.Second)
	go func() {
		var counter = 0
		for range ticker.C {
			log.Println("counter---", counter)
			var details = WriteResponse{
				Message: map[string]interface{}{
					"msg": fmt.Sprintf("%s%d", "check:", counter),
				},
				Sender: SenderOptions{
					AccountId: h.Id,
					UserId:    client.Id,
				},
			}
			h.Dispatch(details)
			counter++
		}
	}()
	client.EventSubscriber[event].Timer = ticker
}
