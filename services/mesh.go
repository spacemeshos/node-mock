package services

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	pb "github.com/spacemeshos/node-mock/spacemesh"
)

// MeshService -
type MeshService struct {
}

// GenesisTime Network genesis time as unix epoch time
func (s MeshService) GenesisTime(ctx context.Context, in *empty.Empty) (*pb.SimpleInt, error) {
	return nil, nil
}

// CurrentLayer Current layer number
func (s MeshService) CurrentLayer(ctx context.Context, in *empty.Empty) (*pb.SimpleInt, error) {
	return nil, nil
}

// CurrentEpoch Current epoch number
func (s MeshService) CurrentEpoch(ctx context.Context, in *empty.Empty) (*pb.SimpleInt, error) {
	return nil, nil
}

// NetId Network ID
func (s MeshService) NetId(ctx context.Context, in *empty.Empty) (*pb.SimpleInt, error) {
	return nil, nil
}

// EpochNumLayers Number of layers per epoch (a network parameter)
func (s MeshService) EpochNumLayers(ctx context.Context, in *empty.Empty) (*pb.SimpleInt, error) {
	return nil, nil
}

////////// Streams

// LayerStream Sent each time layer data changes. Designed for heavy-duty clients. Layer with blocks and transactions.
func (s MeshService) LayerStream(req *empty.Empty, srv pb.MeshService_LayerStreamServer) error {
	return nil
}

// RegisterService registers the grpc service.
func (s MeshService) RegisterService(server *grpc.Server) {
	pb.RegisterMeshServiceServer(server, s)

	// SubscribeOnNewConnections reflection service on gRPC server
	//reflection.Register(server)
}
