package main

import (
	"github.com/kulycloud/api-server/communication"
	"github.com/kulycloud/api-server/config"
	"github.com/kulycloud/api-server/server"
	commonCommunication "github.com/kulycloud/common/communication"
	"github.com/kulycloud/common/logging"
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

	go httpListener(listener.Storage)

	handler := communication.NewApiServerHandler()
	handler.Register(listener)

	if err = listener.Serve(); err != nil {
		logger.Panicw("error serving listener", "error", err)
	}
}

func httpListener(storage *commonCommunication.StorageCommunicator) {
	srv := server.NewServer(storage)
	if err := srv.Start(); err != nil {
		logger.Panicw("error serving http", "error", err)
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

func register() (*commonCommunication.ControlPlaneCommunicator, error) {
	comm := commonCommunication.NewControlPlaneCommunicator()
	err := comm.Connect(config.GlobalConfig.ControlPlaneHost, config.GlobalConfig.ControlPlanePort)
	if err != nil {
		logger.Errorw("Could not connect to control-plane", "error", err)
		return nil, err
	}
	err = comm.RegisterThisService(context.Background(), "api-server", config.GlobalConfig.Host, config.GlobalConfig.Port)
	if err != nil {
		logger.Errorw("Could not register service", "error", err)
		return nil, err
	}
	logger.Info("Registered to control-plane")
	return comm, nil
}
