package services

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/exp/errors/fmt"
	"google.golang.org/grpc"

	"github.com/spacemeshos/node-mock/spacemesh"
)

// NodeService -
type NodeService struct{}

// Echo returns the response for an echo api request
func (s NodeService) Echo(ctx context.Context, in *spacemesh.SimpleString) (*spacemesh.SimpleString, error) {
	return &spacemesh.SimpleString{Value: in.Value}, nil
}

// Version returns the version of the node software as a semver string
func (s NodeService) Version(ctx context.Context, in *empty.Empty) (*spacemesh.SimpleString, error) {
	return &spacemesh.SimpleString{Value: Config.Version}, nil
}

// Build returns the github tag or branch used to build the node
func (s NodeService) Build(ctx context.Context, in *empty.Empty) (*spacemesh.SimpleString, error) {
	return &spacemesh.SimpleString{Value: Config.Build}, nil
}

// Status current node status
func (s NodeService) Status(ctx context.Context, in *empty.Empty) (*spacemesh.NodeStatus, error) {
	return &nodeStatus, nil
}

// SyncStart request that the node start syncing the mesh
func (s NodeService) SyncStart(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	if len(layers) == 0 {
		fmt.Printf("NodeService.SyncStart\n")

		syncStatusBus.Send(
			spacemesh.NodeSyncStatus{
				Status: spacemesh.NodeSyncStatus_SYNCING,
			},
		)

		go startLoadProducer()
	}

	return &empty.Empty{}, nil
}

// SyncStatusStream sync status events
func (s NodeService) SyncStatusStream(empty *empty.Empty, server spacemesh.NodeService_SyncStatusStreamServer) (err error) {
	syncStatusChan, cookie := syncStatusBus.Register()
	defer syncStatusBus.Delete(cookie)

	for {
		select {
		case msg := <-syncStatusChan:
			syncStatus := msg.(spacemesh.NodeSyncStatus)

			err = server.Send(&syncStatus)
			if err != nil {
				fmt.Printf("SyncStatusStream(ERROR): %v\n", err)

				return
			}

			fmt.Printf("SyncStatusStream(OK): %v\n", syncStatus)

		}
	}
}

// ErrorStream node error events
func (s NodeService) ErrorStream(empty *empty.Empty, server spacemesh.NodeService_ErrorStreamServer) (err error) {
	prevError := nodeError

	for {
		if prevError.Type != nodeError.Type {
			err = server.Send(&nodeError)
			if err != nil {
				fmt.Printf("ErrorStream(ERROR): %v\n", err)

				return
			}

			fmt.Printf("ErrorStream(OK): %v\n", nodeError)

			prevError = nodeError
		}

		time.Sleep(1 * time.Second)
	}
}

// InitNode -
func InitNode(s *grpc.Server) {
	spacemesh.RegisterNodeServiceServer(s, NodeService{})
}
