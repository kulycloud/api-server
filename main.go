package main

import (
	"github.com/kulycloud/common/logging"
	"github.com/kulycloud/route-processor/communication"
	"github.com/kulycloud/route-processor/config"
)

func main() {
	initLogger := logging.GetForComponent("init")
	defer logging.Sync()

	err := config.ParseConfig()
	if err != nil {
		initLogger.Fatalw("Error parsing config", "error", err)
	}
	initLogger.Infow("Finished parsing config", "config", config.GlobalConfig)

	comm := communication.NewCommunicator()
	go func() {
		err := comm.Connect()
		if err != nil {
			initLogger.Errorw("Could not connect to control-plane", "error", err)
		}
	}()

	initLogger.Info("Starting listener")
	listener := communication.NewListener()
	err = listener.Start()
	if err != nil {
		initLogger.Panicw("error initializing listener", "error", err)
	}
}