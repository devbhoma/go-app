package main

import (
	coreadmin "goapp/cmd/admin"
	"goapp/cmd/apiserver"
	"goapp/cmd/cli"
	"goapp/cmd/wsserver"
	appconfig "goapp/config"
	"log/slog"
	"os"
)

func main() {
	if err := run(); err != nil {
		slog.Error("cli service stopped", "err", err)
		os.Exit(1)
	}
}

func run() error {
	mainCLI := cli.New()
	cnf := appconfig.Get()

	coreadmin.Boot(mainCLI, cnf)
	apiserver.Boot(mainCLI, cnf) // api server
	wsserver.Boot(mainCLI, cnf)  // socket server

	return mainCLI.Execute()
}
