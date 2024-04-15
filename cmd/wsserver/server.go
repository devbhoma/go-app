package wsserver

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"goapp/cmd/cli"
	socketroutes "goapp/cmd/wsserver/routes"
	appconfig "goapp/config"
	"goapp/internal/authorization"
	"goapp/internal/httpserver"
	"goapp/internal/store"
)

type Server struct {
	Store         store.Base
	Config        appconfig.Config
	Reader        httpserver.Reader
	Authorization authorization.Authorization
}

func Boot(cli *cli.Base, cnf appconfig.Config) {
	cmd := &cobra.Command{
		Use:   "socket",
		Short: "Start the socket server",
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
			engine.LoadHTMLGlob("./cmd/wsserver/templates/*")

			socketroutes.DefineEngine(cnf, dbStore, engine, s.Authorization.Authenticate())

			socketroutes.DefineRoutesV1(cnf, dbStore, engine.Group("api/v1"), s.Authorization.Authenticate())

			err := s.start()
			if err != nil {
				return errors.New("failed to start socket server")
			}
			return nil
		},
	}
	cli.AddCommand(cmd)
}

func (s *Server) start() error {
	s.Config.Port = "2121" // tmp: statics port
	return s.Reader.Run(s.Config.Port)
}
