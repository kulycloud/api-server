package main

import (
	commonCommunication "github.com/kulycloud/common/communication"
	"github.com/kulycloud/common/logging"
	"github.com/kulycloud/route-processor/communication"
	"github.com/kulycloud/route-processor/config"
	"golang.org/x/net/context"
	"time"
)

var logger = logging.GetForComponent("init")

func main() {
	defer logging.Sync()

	err := config.ParseConfig()
	if err != nil {
		logger.Fatalw("Error parsing config", "error", err)
	}
	logger.Infow("Finished parsing config", "config", config.GlobalConfig)

	go registerLoop()

	logger.Info("Starting listener")
	listener := commonCommunication.NewListener(logging.GetForComponent("listener"))
	if err = listener.Setup(config.GlobalConfig.Port); err != nil {
		logger.Panicw("error initializing listener", "error", err)
	}

	handler := communication.NewRouteProcessorHandler()
	handler.Register(listener)

	if err = listener.Serve(); err != nil {
		logger.Panicw("error serving listener", "error", err)
	}
}

func registerLoop() {
	for {
		_, err := register()
		if err == nil {
			break
		}

		logger.Info("Retrying in 5s...")
		time.Sleep(5*time.Second)
	}
}

func register() (*commonCommunication.Communicator, error) {
	comm := commonCommunication.NewCommunicator()
	err := comm.Connect(config.GlobalConfig.ControlPlaneHost, config.GlobalConfig.ControlPlanePort)
	if err != nil {
		logger.Errorw("Could not connect to control-plane", "error", err)
		return nil, err
	}
	err = comm.RegisterThisService(context.Background(), "route-processor", config.GlobalConfig.Host, config.GlobalConfig.Port)
	if err != nil {
		logger.Errorw("Could not register service", "error", err)
		return nil, err
	}
	logger.Info("Registered to control-plane")
	return comm, nil
}
