package main

import (
	"github.com/kulycloud/api-server/communication"
	"github.com/kulycloud/api-server/config"
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

	communication.RegisterToControlPlane()
}
