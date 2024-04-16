package handlers

import (
	"github.com/gin-gonic/gin"
	appconfig "goapp/config"
	"goapp/internal/store"
	"net/http"
)

type Base struct {
	Config appconfig.Config
	Store  store.Base
}

type Handler interface {
	Chat(c *gin.Context)
}

func NewSite(cnf appconfig.Config, str *store.Base, router *gin.RouterGroup) Handler {
	base := &Base{
		Config: cnf,
		Store:  *str,
	}

	router.GET("/", base.Chat)

	return base
}

func (base *Base) Chat(c *gin.Context) {
	c.HTML(http.StatusOK, "chat.gohtml", gin.H{
		"title": "Ws-Chat!",
	})
}
