package sockethandler

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Sender struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	Event      string `json:"event"`
	Data       string `json:"data"`
	ClientId   string `json:"client_id"`
	ClientName string `json:"client_name"`
	Sender     Sender `json:"sender"`
}

type ClientHelper interface {
	Get() *Client
	GetId() string
	GetName() string
	Read(msg *Message) error
	Write(msg Message) error
}

func NewClient(cn *websocket.Conn, id string, clientName string) ClientHelper {
	return &Client{
		Conn: cn,
		Id:   id,
		Name: clientName,
	}
}

func (c *Client) Get() *Client {
	return c
}

func (c *Client) GetId() string {
	return c.Id
}

func (c *Client) GetName() string {
	return c.Name
}

func (c *Client) Read(msg *Message) error {
	return c.Conn.ReadJSON(&msg)
}

func (c *Client) Write(msg Message) error {
	return c.Conn.WriteJSON(msg)
}

func (c *Client) Message(event string, msg string, sender Sender) *Message {
	if sender.Id == "" {
		sender.Id = c.GetId()
		sender.Name = c.GetName()
	}
	return &Message{
		Event:      event,
		Data:       msg,
		ClientId:   c.GetId(),
		ClientName: c.GetName(),
		Sender:     sender,
	}
}
