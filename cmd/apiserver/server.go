package apiserver

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	apiroutes "goapp/cmd/apiserver/routes"
	"goapp/cmd/cli"
	appconfig "goapp/config"
	"goapp/internal/authorization"
	"goapp/internal/httpserver"
	"goapp/internal/store"
	"net/http"
	"time"
)

type Server struct {
	Store         store.Base
	Config        appconfig.Config
	Reader        httpserver.Reader
	Authorization authorization.Authorization
}

func Boot(cli *cli.Base, cnf appconfig.Config) {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "Start the api server",
		RunE: func(cmd *cobra.Command, cliArgs []string) error {
			cmd.SilenceUsage = true
			if cnf.Port == "" {
				return errors.New("port is required")
			}

			dbStore := store.NewStore(cnf)

			s := &Server{
				Config:        cnf,
				Reader:        httpserver.NewServer(cnf.Env),
				Authorization: authorization.New(cnf, dbStore),
			}

			engine := s.Reader.GetEngine()
			engine.LoadHTMLGlob("./cmd/apiserver/templates/*")

			engine.GET("/", func(ctx *gin.Context) {
				ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
					"title": "Hello World",
					"time":  time.Now().Format(time.RFC850),
				})
			})

			apiroutes.DefineRoutes(cnf, dbStore, engine.Group("api/v1"), s.Authorization.Authenticate())

			err := s.start()
			if err != nil {
				return errors.New("failed to start api server")
			}
			return nil
		},
	}
	cli.AddCommand(cmd)
}

func (s *Server) start() error {
	return s.Reader.Run(s.Config.Port)
}
