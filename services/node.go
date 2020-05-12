package services

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/exp/errors/fmt"

	"google.golang.org/grpc"

	pb "github.com/spacemeshos/node-mock/spacemesh"
)

const mockVersion = "0.0.1"
const mockBuild = "180"

var mockStatus = pb.NodeStatus{
	KnownPeers:    5,
	MinPeers:      1,
	MaxPeers:      100,
	IsSynced:      true,
	SyncedLayer:   100,
	CurrentLayer:  100,
	VerifiedLayer: 100,
}

// NodeService -
type NodeService struct {
}

// Echo returns the response for an echo api request
func (s NodeService) Echo(ctx context.Context, in *pb.SimpleString) (*pb.SimpleString, error) {
	return &pb.SimpleString{Value: in.Value}, nil
}

// Version returns the version of the node software as a semver string
func (s NodeService) Version(ctx context.Context, in *empty.Empty) (*pb.SimpleString, error) {
	return &pb.SimpleString{Value: mockVersion}, nil
}

// Build returns the github tag or branch used to build the node
func (s NodeService) Build(ctx context.Context, in *empty.Empty) (*pb.SimpleString, error) {
	return &pb.SimpleString{Value: mockBuild}, nil
}

// Status current node status
func (s NodeService) Status(context.Context, *empty.Empty) (*pb.NodeStatus, error) {
	return &mockStatus, nil
}

// SyncStart request that the node start syncing the mesh
func (s NodeService) SyncStart(context.Context, *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

// Setream API =====

func statusLoader(svr pb.NodeService_SyncStatusStreamServer) error {
	var err error

	for {
		err = svr.Send(&pb.NodeSyncStatus{})
		if err != nil {
			fmt.Printf("Send(ERROR): %v\n", err)

			break
		}

		fmt.Printf("Send(OK).\n")

		time.Sleep(10 * time.Second)
	}

	return err
}

// SyncStatusStream sync status events
func (s NodeService) SyncStatusStream(empty *empty.Empty, svr pb.NodeService_SyncStatusStreamServer) error {
	var err error

	fmt.Printf("SyncStatusStream: %v\n", svr)

	statusLoader(svr)

	return err
}

// ErrorStream node error events
func (s NodeService) ErrorStream(empty *empty.Empty, svr pb.NodeService_ErrorStreamServer) error {
	return nil
}

// RegisterService registers the grpc service.
func (s NodeService) RegisterService(server *grpc.Server) {
	pb.RegisterNodeServiceServer(server, s)

	// SubscribeOnNewConnections reflection service on gRPC server
	//reflection.Register(server)
}
