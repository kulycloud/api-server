package communication

import (
	"context"
	"fmt"
	"github.com/kulycloud/common/logging"
	"github.com/kulycloud/protocol/common"
	protoRouteProcessor "github.com/kulycloud/protocol/route-processor"
	"github.com/kulycloud/route-processor/config"
	"google.golang.org/grpc"
	"net"
)


var _ protoRouteProcessor.RouteProcessorServer = &Listener{}

var logger = logging.GetForComponent("communication")

type Listener struct {
	server *grpc.Server
	listener net.Listener
}

func NewListener() *Listener {
	return &Listener{}
}

func (listener *Listener) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", config.GlobalConfig.Port))
	if err != nil {
		return err
	}
	listener.listener = lis
	listener.server = grpc.NewServer()
	protoRouteProcessor.RegisterRouteProcessorServer(listener.server, listener)
	logger.Infow("serving", "port", config.GlobalConfig.Port)
	return listener.server.Serve(listener.listener)
}

func (listener *Listener) ProcessRoute(ctx context.Context, request *protoRouteProcessor.RouteProcessorRequest) (*protoRouteProcessor.RouteProcessorResponse, error) {
	return &protoRouteProcessor.RouteProcessorResponse{Status: protoRouteProcessor.RouteProcessorStatus_OK}, nil
}

func (listener *Listener) Ping(ctx context.Context, empty *common.Empty) (*common.Empty, error) {
	return empty, nil
}
