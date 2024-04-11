package socketroutes

import (
	"github.com/gin-gonic/gin"
	sockethandler "goapp/cmd/wsserver/routes/handlers/socket"
	appconfig "goapp/config"
	"goapp/internal/store"
)

func DefineRoutes(cnf appconfig.Config, str *store.Base, v1 *gin.RouterGroup, authHandler gin.HandlerFunc) {

	socket := sockethandler.New(cnf, str, v1)

	v1.Use(authHandler)
	{
		socket.PrivateRoutes(v1)
	}
}
