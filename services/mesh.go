package services

import (
	"context"

	"golang.org/x/exp/errors/fmt"
	"google.golang.org/grpc"

	v1 "github.com/spacemeshos/api/release/go/spacemesh/v1"
)

// MeshService -
type MeshService struct{}

// GenesisTime Network genesis time as unix epoch time
func (s MeshService) GenesisTime(ctx context.Context, request *v1.GenesisTimeRequest) (*v1.GenesisTimeResponse, error) {
	return &v1.GenesisTimeResponse{
		Unixtime: &v1.SimpleInt{Value: network.GenesisTime},
	}, nil
}

// CurrentLayer Current layer number
func (s MeshService) CurrentLayer(ctx context.Context, request *v1.CurrentLayerRequest) (*v1.CurrentLayerResponse, error) {
	return &v1.CurrentLayerResponse{
		Layernum: &v1.SimpleInt{Value: network.GetCurrentLayerNumber()},
	}, nil
}

// CurrentEpoch Current epoch number
func (s MeshService) CurrentEpoch(ctx context.Context, request *v1.CurrentEpochRequest) (*v1.CurrentEpochResponse, error) {
	return &v1.CurrentEpochResponse{
		Epochnum: &v1.SimpleInt{Value: network.GetCurrentEpochNumber()},
	}, nil
}

// NetId Network ID
func (s MeshService) NetId(ctx context.Context, request *v1.NetIdRequest) (*v1.NetIdResponse, error) {
	return &v1.NetIdResponse{
		Netid: &v1.SimpleInt{Value: network.NetId},
	}, nil
}

// EpochNumLayers Number of layers per epoch (a network parameter)
func (s MeshService) EpochNumLayers(ctx context.Context, request *v1.EpochNumLayersRequest) (*v1.EpochNumLayersResponse, error) {
	return &v1.EpochNumLayersResponse{
		Numlayers: &v1.SimpleInt{Value: network.EpochNumLayers},
	}, nil
}

// LayerDuration Layer duration (a network parameter)
func (s MeshService) LayerDuration(ctx context.Context, request *v1.LayerDurationRequest) (*v1.LayerDurationResponse, error) {
	return &v1.LayerDurationResponse{
		Duration: &v1.SimpleInt{Value: network.LayerDuration},
	}, nil
}

// MaxTransactionsPerSecond Number of transactions per second (a network parameter)
func (s MeshService) MaxTransactionsPerSecond(ctx context.Context, request *v1.MaxTransactionsPerSecondRequest) (*v1.MaxTransactionsPerSecondResponse, error) {
	return &v1.MaxTransactionsPerSecondResponse{
		Maxtxpersecond: &v1.SimpleInt{Value: network.MaxTransactionsPerSecond},
	}, nil
}

// AccountMeshDataQuery Get account data query
func (s MeshService) AccountMeshDataQuery(ctx context.Context, request *v1.AccountMeshDataQueryRequest) (*v1.AccountMeshDataQueryResponse, error) {
	return &v1.AccountMeshDataQueryResponse{}, nil
}

// LayersQuery Layers data query
func (s MeshService) LayersQuery(ctx context.Context, request *v1.LayersQueryRequest) (result *v1.LayersQueryResponse, err error) {
	result = new(v1.LayersQueryResponse)

	for i := request.StartLayer; i <= request.EndLayer; i++ {
		result.Layer = append(result.Layer, network.Layers[i].Export())
	}

	return
}

// Setream API =====

// AccountMeshDataStream A stream of transactions and activations from an account.
// Includes simple coin transactions with the account as the destination.
func (s MeshService) AccountMeshDataStream(request *v1.AccountMeshDataStreamRequest, server v1.MeshService_AccountMeshDataStreamServer) error {
	return nil
}

// LayerStream Sent each time layer data changes. Designed for heavy-duty clients. Layer with blocks and transactions.
func (s MeshService) LayerStream(request *v1.LayerStreamRequest, server v1.MeshService_LayerStreamServer) (err error) {
	layerChan, cookie := layerBus.Register()
	defer layerBus.Delete(cookie)

	fmt.Printf("LayerStream: started\n")

	for {
		select {
		case msg := <-layerChan:
			fmt.Printf("LayerStream: request\n")
			layer := msg.(*v1.Layer)

			response := &v1.LayerStreamResponse{
				Layer: layer,
			}

			err = server.Send(response)
			if err != nil {
				fmt.Printf("LayerStream(ERROR): %v\n", err)

				return
			}

			fmt.Printf("LayerStream(OK): %d - %s\n", response.Layer.GetNumber(), response.Layer.GetStatus().String())
		}
	}
}

// InitMesh -
func InitMesh(s *grpc.Server) {
	v1.RegisterMeshServiceServer(s, MeshService{})
}
