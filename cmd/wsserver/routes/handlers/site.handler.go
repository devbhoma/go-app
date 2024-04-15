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
	Page(c *gin.Context)
}

func NewSite(cnf appconfig.Config, str *store.Base, router *gin.RouterGroup) Handler {
	base := &Base{
		Config: cnf,
		Store:  *str,
	}

	router.GET("/", base.Page)
	router.GET("/chat", base.Chat)

	return base
}

func (base *Base) Page(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "WS!",
	})
}

func (base *Base) Chat(c *gin.Context) {
	c.HTML(http.StatusOK, "site.gohtml", gin.H{
		"title": "Ws-Chat!",
	})
}
