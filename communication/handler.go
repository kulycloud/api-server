package communication

import (
	commonCommunication "github.com/kulycloud/common/communication"
	"github.com/kulycloud/common/logging"
	protoApiServer "github.com/kulycloud/protocol/api-server"
)


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
