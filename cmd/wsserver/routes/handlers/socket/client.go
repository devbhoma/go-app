package sockethandler

import (
	"github.com/gorilla/websocket"
	"time"
)

type Subscriber struct {
	Timer *time.Ticker
}

type Client struct {
	WsConn          *websocket.Conn
	Id              string
	EventSubscriber map[string]*Subscriber
	//DataProducer    pulsar_client.ClientProducer
	//DataConsumer    pulsar_client.ClientConsumer
}

type ClientHelper interface {
	Get() *Client
	RemoveTimer(event string)
}

func NewClient(cn *websocket.Conn, id string, account_id string) ClientHelper {
	return &Client{
		WsConn:          cn,
		Id:              id,
		EventSubscriber: make(map[string]*Subscriber),
	}
}

func (c *Client) Get() *Client {
	return c
}

func (c *Client) RemoveTimer(event string) {
	if c.EventSubscriber[event] != nil {
		subscriber := c.EventSubscriber[event]

		if subscriber.Timer != nil {
			subscriber.Timer.Stop()
			subscriber.Timer = nil

			delete(c.EventSubscriber, event)
		}
	}
}
