package communication

import (
	"context"
	"fmt"
	protoCommon "github.com/kulycloud/protocol/common"
	protoControlPlane "github.com/kulycloud/protocol/control-plane"
	"github.com/kulycloud/route-processor/config"
	"google.golang.org/grpc"
)

type Communicator struct {
	connection *grpc.ClientConn
	controlPlaneClient protoControlPlane.ControlPlaneClient
}

func NewCommunicator() *Communicator {
	return &Communicator{}
}

func (communicator *Communicator) Connect() error {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%v", config.GlobalConfig.ControlPlaneHost, config.GlobalConfig.ControlPlanePort), grpc.WithInsecure())
	if err != nil {
		return err
	}
	communicator.controlPlaneClient = protoControlPlane.NewControlPlaneClient(conn)
	_, err = communicator.controlPlaneClient.RegisterComponent(context.TODO(), &protoControlPlane.RegisterComponentRequest {
		Type:     "route-processor",
		Endpoint: &protoCommon.Endpoint{
			Host: config.GlobalConfig.Host,
			Port: config.GlobalConfig.Port,
		},
	})

	if err != nil {
		return fmt.Errorf("error from control-plane during connection: %w", err)
	}
	logger.Info("connected to control-plane")
	return nil
}
