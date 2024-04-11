package sockethandler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gorilla/websocket"
	appconfig "goapp/config"
	"goapp/internal/store"
	"net/http"
	"time"
)

type ClientSubscriber struct {
	UserId    int
	AccountId int
	Timer     *time.Ticker
}

type Base struct {
	Upgrade    websocket.Upgrader
	WsConnHubs map[string]*Hub
	Subscriber map[string]*ClientSubscriber
}

type Handler interface {
	PrivateRoutes(router *gin.RouterGroup)
	Broadcast(ctx *gin.Context)
	Client(ctx *gin.Context)
	WebSocket(ctx *gin.Context)
	Subscribe(ctx *gin.Context)
	Unsubscribe(ctx *gin.Context)
}

func New(cnf appconfig.Config, str *store.Base, router *gin.RouterGroup) Handler {

	base := &Base{
		Upgrade: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		WsConnHubs: make(map[string]*Hub),
		Subscriber: make(map[string]*ClientSubscriber),
	}

	router.GET("/ws/websocket", base.WebSocket)
	router.POST("/ws/broadcast", base.Broadcast)
	router.POST("/ws/emit", base.Client)
	router.GET("/ws/subscribe/:key", base.Subscribe)
	router.GET("/ws/unsubscribe/:key", base.Unsubscribe)

	return base
}

func (b *Base) PrivateRoutes(router *gin.RouterGroup) {
	router.GET("/ping", b.Ping)
}

func (b *Base) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (b *Base) WebSocket(ctx *gin.Context) {
	accountId := "6"
	userId := "2"

	wsConn, err := b.Upgrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Printf("error in upgrading connection", "err", err)
		return
	}

	newClient := NewClient(wsConn, userId, accountId)
	client := newClient.Get()

	wsHub := b.WsConnHubs[accountId]
	if wsHub == nil {
		newHub := NewConnHub(accountId, userId)
		wsHub = newHub.Get()

		b.WsConnHubs[accountId] = wsHub
		go wsHub.Write()
	}
	defer func() {
		wsHub.RemoveClient(client, "")
		if err0 := wsConn.Close(); err0 != nil {
			fmt.Println("Error while closing socket!!!")
			fmt.Printf("error in closing socket", "err", err0)
			return
		}
	}()

	wsHub.Clients[userId] = client
	wsHub.Read(client)
}

func (b *Base) Broadcast(ctx *gin.Context) {
	accountId := "6"
	var req map[string]interface{}
	err := ctx.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		fmt.Println("binding error:", err)
		ctx.JSON(http.StatusOK, gin.H{
			"status":  false,
			"message": "error in parsing request",
		})
		return
	}

	var hub = b.WsConnHubs[accountId]
	if hub != nil {
		var details = WriteResponse{
			Message: req,
			Sender: SenderOptions{
				AccountId: accountId,
			},
		}
		hub.Broadcast <- details

	}
}

func (b *Base) Client(ctx *gin.Context) {
	accountId := "6"
	userId := "2"

	var req map[string]interface{}
	_ = ctx.BindJSON(&req)

	var hub = b.WsConnHubs[accountId]
	if hub != nil {
		var details = WriteResponse{
			Message: req,
			Sender: SenderOptions{
				AccountId: accountId,
				UserId:    userId,
			},
		}
		hub.SendToWrite <- details
	}
}

func (b *Base) Unsubscribe(ctx *gin.Context) {
	eventKey := ctx.Param("key")

	accountId := "6"
	userId := "2"

	var wsHub = b.WsConnHubs[accountId]
	if wsHub != nil {
		client := wsHub.Clients[userId]
		client.RemoveTimer(eventKey)

		ctx.JSON(http.StatusOK, gin.H{
			"status":  true,
			"event":   eventKey,
			"type":    "unsubscribe",
			"message": "this event unsubscribed!",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  false,
			"message": "Invalid client ws connection",
		})
	}
}

func (b *Base) Subscribe(ctx *gin.Context) {
	EventKey := ctx.Param("key")

	accountId := "6"
	userId := "2"

	var wsHub = b.WsConnHubs[accountId]
	if wsHub != nil {
		client := wsHub.Clients[userId]
		wsHub.TestSubscribeLogs(client, EventKey)
		ctx.JSON(http.StatusOK, gin.H{
			"status":  true,
			"event":   EventKey,
			"type":    "subscribe",
			"message": "this event subscribed!",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  false,
			"message": "Invalid client ws connection",
		})
	}
}
