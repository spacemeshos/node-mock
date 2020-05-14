package services

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/exp/errors/fmt"
	"google.golang.org/grpc"

	"github.com/spacemeshos/node-mock/spacemesh"
)

const mockVersion = "0.0.1"
const mockBuild = "1"

var syncStatusMock = []spacemesh.NodeSyncStatus{
	{Status: spacemesh.NodeSyncStatus_SYNCING},
	{Status: spacemesh.NodeSyncStatus_SYNCED},
	{Status: spacemesh.NodeSyncStatus_NEW_LAYER_VERIFIED},
	{Status: spacemesh.NodeSyncStatus_NEW_TOP_LAYER},
	{Status: spacemesh.NodeSyncStatus_NEW_LAYER_VERIFIED},
}

var syncPosition int
var syncStatus spacemesh.NodeSyncStatus
var nodeStatus = spacemesh.NodeStatus{
	KnownPeers:    50,
	MinPeers:      1,
	MaxPeers:      10,
	IsSynced:      true,
	SyncedLayer:   100,
	CurrentLayer:  100,
	VerifiedLayer: 90,
}

var nodeError spacemesh.NodeError

// NodeService -
type NodeService struct{}

// Echo returns the response for an echo api request
func (s NodeService) Echo(ctx context.Context, in *spacemesh.SimpleString) (*spacemesh.SimpleString, error) {
	return &spacemesh.SimpleString{Value: in.Value}, nil
}

// Version returns the version of the node software as a semver string
func (s NodeService) Version(ctx context.Context, in *empty.Empty) (*spacemesh.SimpleString, error) {
	return &spacemesh.SimpleString{Value: mockVersion}, nil
}

// Build returns the github tag or branch used to build the node
func (s NodeService) Build(ctx context.Context, in *empty.Empty) (*spacemesh.SimpleString, error) {
	return &spacemesh.SimpleString{Value: mockBuild}, nil
}

// Status current node status
func (s NodeService) Status(ctx context.Context, in *empty.Empty) (*spacemesh.NodeStatus, error) {
	return &nodeStatus, nil
}

// SyncStart request that the node start syncing the mesh
func (s NodeService) SyncStart(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

// SyncStatusStream sync status events
func (s NodeService) SyncStatusStream(empty *empty.Empty, server spacemesh.NodeService_SyncStatusStreamServer) (err error) {
	prevStatus := syncStatus

	for {
		if prevStatus.Status != syncStatus.Status {
			err = server.Send(&syncStatus)
			if err != nil {
				fmt.Printf("SyncStatusStream(ERROR): %v\n", err)

				return
			}

			fmt.Printf("SyncStatusStream(OK): %v\n", syncStatus)

			prevStatus = syncStatus
		}

		time.Sleep(1 * time.Second)
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

			fmt.Printf("ErrorStream(OK): %v\n", syncStatus)

			prevError = nodeError
		}

		time.Sleep(1 * time.Second)
	}
}

func updateSyncStatus() {
	if syncPosition >= len(syncStatusMock) {
		syncPosition = len(syncStatusMock) - 1
	}

	syncStatus = syncStatusMock[syncPosition]

	syncPosition++
}

func statusLoadProducer() {
	for {
		updateSyncStatus()

		fmt.Printf("statusLoadProducer: %s\n", syncStatus.Status.String())

		time.Sleep(10 * time.Second)
	}
}

// InitNode -
func InitNode(s *grpc.Server) {
	go statusLoadProducer()

	spacemesh.RegisterNodeServiceServer(s, NodeService{})
}
