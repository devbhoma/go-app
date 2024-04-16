package socketroutes

import (
	"github.com/gin-gonic/gin"
	"goapp/cmd/wsserver/routes/handlers"
	sockethandler "goapp/cmd/wsserver/routes/handlers/socket"
	appconfig "goapp/config"
	"goapp/internal/store"
)

func DefineEngine(cnf appconfig.Config, dbStore *store.Base, engine *gin.Engine, authHandler gin.HandlerFunc) {

	handlers.NewSite(cnf, dbStore, engine.Group("/"))

}

func DefineRoutesV1(cnf appconfig.Config, str *store.Base, v1 *gin.RouterGroup, authHandler gin.HandlerFunc) {

	sockethandler.New(cnf, str, v1)

	v1.Use(authHandler)
	{
		//socket.PrivateRoutes(v1)
	}
}
