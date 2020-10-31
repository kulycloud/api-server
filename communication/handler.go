package communication

import (
	"context"
	commonCommunication "github.com/kulycloud/common/communication"
	"github.com/kulycloud/common/logging"
	protoRouteProcessor "github.com/kulycloud/protocol/route-processor"
)


var _ protoRouteProcessor.RouteProcessorServer = &RouteProcessorHandler{}

var logger = logging.GetForComponent("handler")

type RouteProcessorHandler struct {
	protoRouteProcessor.UnimplementedRouteProcessorServer
}

func NewRouteProcessorHandler() *RouteProcessorHandler {
	return &RouteProcessorHandler{}
}

func (handler *RouteProcessorHandler) Register(listener *commonCommunication.Listener) {
	protoRouteProcessor.RegisterRouteProcessorServer(listener.Server, handler)
}

func (handler *RouteProcessorHandler) ProcessRoute(ctx context.Context, request *protoRouteProcessor.RouteProcessorRequest) (*protoRouteProcessor.RouteProcessorResponse, error) {
	logger.Infow("processing route", "route", request.Data)
	return &protoRouteProcessor.RouteProcessorResponse{Status: protoRouteProcessor.RouteProcessorStatus_OK}, nil
}

