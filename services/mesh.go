package services

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/exp/errors/fmt"
	"google.golang.org/grpc"

	"github.com/spacemeshos/node-mock/spacemesh"
)

// MeshService -
type MeshService struct{}

// GenesisTime Network genesis time as unix epoch time
func (s MeshService) GenesisTime(ctx context.Context, in *empty.Empty) (*spacemesh.SimpleInt, error) {
	return &spacemesh.SimpleInt{Value: uint64(genesisTime.Unix())}, nil
}

// CurrentLayer Current layer number
func (s MeshService) CurrentLayer(ctx context.Context, in *empty.Empty) (*spacemesh.SimpleInt, error) {
	return &spacemesh.SimpleInt{Value: uint64(nodeStatus.CurrentLayer)}, nil
}

// CurrentEpoch Current epoch number
func (s MeshService) CurrentEpoch(ctx context.Context, in *empty.Empty) (*spacemesh.SimpleInt, error) {
	return &spacemesh.SimpleInt{Value: currentEpoch}, nil
}

// NetId Network ID
func (s MeshService) NetId(ctx context.Context, in *empty.Empty) (*spacemesh.SimpleInt, error) {
	return &spacemesh.SimpleInt{Value: Config.NetID}, nil
}

// EpochNumLayers Number of layers per epoch (a network parameter)
func (s MeshService) EpochNumLayers(ctx context.Context, in *empty.Empty) (*spacemesh.SimpleInt, error) {
	return &spacemesh.SimpleInt{Value: Config.Layers.PerEpoch}, nil
}

// Setream API =====

// LayerStream Sent each time layer data changes. Designed for heavy-duty clients. Layer with blocks and transactions.
func (s MeshService) LayerStream(req *empty.Empty, server spacemesh.MeshService_LayerStreamServer) (err error) {
	layerChan, cookie := layerBus.Register()
	defer layerBus.Delete(cookie)

	for {
		select {
		case msg := <-layerChan:
			layer := msg.(spacemesh.Layer)

			err = server.Send(&layer)
			if err != nil {
				fmt.Printf("LayerStream(ERROR): %v\n", err)

				return
			}

			fmt.Printf("LayerStream(OK): %d - %s\n", layer.GetNumber(), layer.GetStatus().String())

		}
	}
}

// InitMesh -
func InitMesh(s *grpc.Server) {
	spacemesh.RegisterMeshServiceServer(s, MeshService{})
}
