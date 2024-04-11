package handlers

import (
	"github.com/gin-gonic/gin"
	appconfig "goapp/config"
	"goapp/internal/store"
	"net/http"
)

type SiteHandler interface {
	PrivateRoutes(router *gin.RouterGroup)
	Ping(ctx *gin.Context)
	User(ctx *gin.Context)
}

type Site struct {
	Store  *store.Base
	Config appconfig.Config
}

func NewSiteHandler(cnf appconfig.Config, str *store.Base, router *gin.RouterGroup) SiteHandler {
	site := &Site{
		Store:  str,
		Config: cnf,
	}

	router.GET("/ping", site.Ping)
	return site
}

func (s *Site) PrivateRoutes(router *gin.RouterGroup) {
	router.GET("/user", s.User)
}

func (s *Site) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"ping": "pong",
	})
}

func (s *Site) User(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"user": true,
	})
}
