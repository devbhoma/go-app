package sockethandler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	appconfig "goapp/config"
	"goapp/internal/httpserver"
	"goapp/internal/store"
	"net/http"
)

type Base struct {
	Upgrade websocket.Upgrader
	Rooms   map[string]*Room
}

type Handler interface {
	WebSocket(ctx *gin.Context)
	Broadcast(ctx *gin.Context)
	Emit(ctx *gin.Context)
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
		Rooms: make(map[string]*Room),
	}

	router.GET("/ws/websocket/:roomId", base.WebSocket)
	router.POST("/ws/broadcast", base.Broadcast)
	router.POST("/ws/emit/:clientId", base.Emit)

	return base
}

func (b *Base) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (b *Base) WebSocket(ctx *gin.Context) {
	roomId := ctx.Param("roomId")
	clientId := httpserver.GetQueryContextString(ctx, "client_id")
	roomName := httpserver.GetQueryContextString(ctx, "room_name")
	clientName := httpserver.GetQueryContextString(ctx, "client_name")

	con, err := b.Upgrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Printf("error in upgrading connection, err: %s", err)
		return
	}
	isNewRoom := false
	isNewClient := false

	room := b.Rooms[roomId]
	if room == nil {
		newRoom := NewRoom(roomId, roomName)

		room = newRoom.Get()
		isNewRoom = true
		b.Rooms[roomId] = room
		go room.Write()
	}

	client := room.Clients[clientId]
	if client == nil {
		newClient := NewClient(con, clientId, clientName)
		client = newClient.Get()
		isNewClient = true
	}

	if isNewRoom {
		room.Emit(Message{
			Event:      "chat:message::event",
			Data:       fmt.Sprintf("%s room created", roomName),
			ClientId:   clientId,
			ClientName: clientName,
			Sender: Sender{
				Id:   clientId,
				Name: clientName,
			},
		})
	}

	if isNewClient {
		room.EmitAll(Message{
			Event:      "chat:message::event",
			Data:       fmt.Sprintf("%s joined", clientName),
			ClientId:   clientId,
			ClientName: clientName,
			Sender: Sender{
				Id:   clientId,
				Name: clientName,
			},
		})
	}

	defer func() {
		room.RemoveClient(client)
		if err := con.Close(); err != nil {
			fmt.Println("Error while closing socket, err:", err)
		}
	}()

	room.Clients[clientId] = client
	room.Read(client)

}

func (b *Base) Broadcast(ctx *gin.Context) {

	//var req map[string]interface{}
	//err := ctx.ShouldBindBodyWith(&req, binding.JSON)
	//if err != nil {
	//	fmt.Println("binding error:", err)
	//	ctx.JSON(http.StatusOK, gin.H{
	//		"status":  false,
	//		"message": "error in parsing request",
	//	})
	//	return
	//}
	//
	//if len(b.WsConnHubs) > 0 {
	//	for clientId := range b.WsConnHubs {
	//		var hub = b.WsConnHubs[clientId]
	//		if hub != nil {
	//			var details = WriteResponse{
	//				Message: req,
	//				Sender: SenderOptions{
	//					Id: clientId,
	//				},
	//			}
	//			hub.Broadcast <- details
	//		}
	//	}
	//}
}

func (b *Base) Emit(ctx *gin.Context) {
	//clientId := ctx.Param("clientId")
	//var req map[string]interface{}
	//_ = ctx.BindJSON(&req)
	//
	//if hub, ok := b.WsConnHubs[clientId]; ok && hub != nil {
	//	details := WriteResponse{
	//		Message: req,
	//		Sender: SenderOptions{
	//			Id: clientId,
	//		},
	//	}
	//	hub.SendToWrite <- details
	//
	//	ctx.JSON(http.StatusOK, gin.H{
	//		"status":  true,
	//		"message": "Message Emitted",
	//	})
	//	return
	//}
	//ctx.JSON(http.StatusOK, gin.H{
	//	"status":  false,
	//	"message": "client not found",
	//})
}
