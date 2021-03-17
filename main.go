package main

import (
	"github.com/kulycloud/api-server/communication"
	"github.com/kulycloud/api-server/config"
	"github.com/kulycloud/api-server/server"
	"github.com/kulycloud/common/logging"
)

var logger = logging.GetForComponent("init")

func main() {
	defer logging.Sync()

	err := config.ParseConfig()
	if err != nil {
		logger.Fatalw("Error parsing config", "error", err)
	}
	logger.Infow("Finished parsing config", "config", config.GlobalConfig)

	serveErrs := communication.RegisterToControlPlane()

	go startHTTPServer()

	err = <-serveErrs
	if err != nil {
		logger.Panicw("error serving listener", "error", err)
	}
}

func startHTTPServer() {
	srv := server.NewServer(communication.ControlPlane.Storage)
	if err := srv.Start(); err != nil {
		logger.Panicw("error serving http", "error", err)
	}
}
