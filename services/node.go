package services

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/exp/errors/fmt"
	"google.golang.org/grpc"

	v1 "github.com/spacemeshos/api/release/go/spacemesh/v1"
)

// NodeService -
type NodeService struct{}

// Echo returns the response for an echo api request
func (s NodeService) Echo(ctx context.Context, request *v1.EchoRequest) (*v1.EchoResponse, error) {
	return &v1.EchoResponse{
		Msg: &v1.SimpleString{Value: request.GetMsg().Value},
	}, nil
}

// Version returns the version of the node software as a semver string
func (s NodeService) Version(ctx context.Context, in *empty.Empty) (*v1.VersionResponse, error) {
	return &v1.VersionResponse{
		VersionString: &v1.SimpleString{Value: Config.Version},
	}, nil
}

// Build returns the github tag or branch used to build the node
func (s NodeService) Build(ctx context.Context, in *empty.Empty) (*v1.BuildResponse, error) {
	return &v1.BuildResponse{
		BuildString: &v1.SimpleString{Value: Config.Build},
	}, nil
}

// Status current node status
func (s NodeService) Status(ctx context.Context, request *v1.StatusRequest) (*v1.StatusResponse, error) {
	return &v1.StatusResponse{}, nil
}

// SyncStart request that the node start syncing the mesh
func (s NodeService) SyncStart(ctx context.Context, request *v1.SyncStartRequest) (*v1.SyncStartResponse, error) {
	if len(layers) == 0 {
		fmt.Printf("NodeService.SyncStart\n")

		/*syncStatusBus.Send(
			v1.NodeSyncStatus{
				Status: v1.NodeSyncStatus_SYNCING,
			},
		)*/

		//go startLoadProducer()
	}

	return &v1.SyncStartResponse{}, nil
}

// Shutdown Request that the node initiate graceful shutdown
func (s NodeService) Shutdown(ctx context.Context, request *v1.ShutdownRequest) (*v1.ShutdownResponse, error) {
	return &v1.ShutdownResponse{}, nil
}

// StatusStream sync status events
func (s NodeService) StatusStream(request *v1.StatusStreamRequest, server v1.NodeService_StatusStreamServer) (err error) {
	/*syncStatusChan, cookie := syncStatusBus.Register()
	defer syncStatusBus.Delete(cookie)

	for {
		select {
		case msg := <-syncStatusChan:
			syncStatus := msg.(v1.NodeSyncStatus)

			err = server.Send(&syncStatus)
			if err != nil {
				fmt.Printf("SyncStatusStream(ERROR): %v\n", err)

				return
			}

			fmt.Printf("SyncStatusStream(OK): %v\n", syncStatus)

		}
	}*/

	return nil
}

// ErrorStream node error events
func (s NodeService) ErrorStream(request *v1.ErrorStreamRequest, server v1.NodeService_ErrorStreamServer) (err error) {
	/*prevError := nodeError

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
	}*/

	return nil
}

// InitNode -
func InitNode(s *grpc.Server) {
	v1.RegisterNodeServiceServer(s, NodeService{})
}
