package apiroutes

import (
	"github.com/gin-gonic/gin"
	"goapp/cmd/apiserver/routes/handlers"
	appconfig "goapp/config"
	"goapp/internal/store"
)

func DefineRoutes(cnf appconfig.Config, str *store.Base, v1 *gin.RouterGroup, authHandler gin.HandlerFunc) {
	handlers.NewAuthHandler(cnf, str, v1)
	site := handlers.NewSiteHandler(cnf, str, v1)

	v1.Use(authHandler)
	{
		site.PrivateRoutes(v1)
	}
}
