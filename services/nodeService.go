package services

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"

	pb "github.com/spacemeshos/node-mock/spacemesh"
)

// NodeService -
type NodeService struct {
}

// Echo returns the response for an echo api request
func (s NodeService) Echo(ctx context.Context, in *pb.SimpleString) (*pb.SimpleString, error) {
	fmt.Printf("ECHO: %s\n", in.GetValue())

	return &pb.SimpleString{Value: in.Value}, nil
}

// Version returns the version of the node software as a semver string
func (s NodeService) Version(ctx context.Context, in *empty.Empty) (*pb.SimpleString, error) {
	return &pb.SimpleString{Value: "version"}, nil
}

// Build returns the github tag or branch used to build the node
func (s NodeService) Build(ctx context.Context, in *empty.Empty) (*pb.SimpleString, error) {
	return &pb.SimpleString{Value: "build"}, nil
}

// Status current node status
func (s NodeService) Status(context.Context, *empty.Empty) (*pb.NodeStatus, error) {
	return &pb.NodeStatus{}, nil
}

// SyncStart request that the node start syncing the mesh
func (s NodeService) SyncStart(context.Context, *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

// SyncStatusStream sync status events
func (s NodeService) SyncStatusStream(*empty.Empty, pb.NodeService_SyncStatusStreamServer) error {
	return nil
}

// ErrorStream node error events
func (s NodeService) ErrorStream(*empty.Empty, pb.NodeService_ErrorStreamServer) error {
	return nil
}

// RegisterService registers the grpc service.
func (s NodeService) RegisterService(server *grpc.Server) {
	pb.RegisterNodeServiceServer(server, s)

	// SubscribeOnNewConnections reflection service on gRPC server
	reflection.Register(server)
}
