package communication

import (
	"github.com/kulycloud/api-server/config"
	"github.com/kulycloud/api-server/server"
	commonCommunication "github.com/kulycloud/common/communication"
	"github.com/kulycloud/common/logging"
	protoApiServer "github.com/kulycloud/protocol/api-server"
)

var ControlPlane *commonCommunication.ControlPlaneCommunicator

var _ protoApiServer.ApiServerServer = &ApiServerHandler{}

var logger = logging.GetForComponent("handler")

type ApiServerHandler struct {
	protoApiServer.UnimplementedApiServerServer
}

func NewApiServerHandler() *ApiServerHandler {
	return &ApiServerHandler{}
}

func (handler *ApiServerHandler) Register(listener *commonCommunication.Listener) {
	protoApiServer.RegisterApiServerServer(listener.Server, handler)
}

func RegisterToControlPlane() {
	communicator := commonCommunication.RegisterToControlPlane("api-server",
		config.GlobalConfig.Host, config.GlobalConfig.Port,
		config.GlobalConfig.ControlPlaneHost, config.GlobalConfig.ControlPlanePort)

	logger.Info("Starting listener")
	listener := commonCommunication.NewListener(logging.GetForComponent("listener"))
	if err := listener.Setup(config.GlobalConfig.Port); err != nil {
		logger.Panicw("error initializing listener", "error", err)
	}

	go httpListener(listener.Storage)

	handler := NewApiServerHandler()
	handler.Register(listener)
	serveErrs := listener.Serve()

	ControlPlane = <-communicator
	err := <-serveErrs
	if err != nil {
		logger.Panicw("error serving listener", "error", err)
	}
}

func httpListener(storage *commonCommunication.StorageCommunicator) {
	srv := server.NewServer(storage)
	if err := srv.Start(); err != nil {
		logger.Panicw("error serving http", "error", err)
	}
}
