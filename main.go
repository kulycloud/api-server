package main

import (
	commonCommunication "github.com/kulycloud/common/communication"
	"github.com/kulycloud/common/logging"
	"github.com/kulycloud/route-processor/communication"
	"github.com/kulycloud/route-processor/config"
	"golang.org/x/net/context"
)

func main() {
	initLogger := logging.GetForComponent("init")
	defer logging.Sync()

	err := config.ParseConfig()
	if err != nil {
		initLogger.Fatalw("Error parsing config", "error", err)
	}
	initLogger.Infow("Finished parsing config", "config", config.GlobalConfig)

	comm := commonCommunication.NewCommunicator()
	go func() {
		err := comm.Connect(config.GlobalConfig.ControlPlaneHost, config.GlobalConfig.ControlPlanePort)
		if err != nil {
			initLogger.Errorw("Could not connect to control-plane", "error", err)
		}
		err = comm.RegisterThisService(context.Background(), "route-processor", config.GlobalConfig.Host, config.GlobalConfig.Port)
		if err != nil {
			initLogger.Errorw("Could not register service", "error", err)
		}
		initLogger.Info("Registered to control-plane")
	}()

	initLogger.Info("Starting listener")
	listener := communication.NewListener()
	err = listener.Start()
	if err != nil {
		initLogger.Panicw("error initializing listener", "error", err)
	}
}
