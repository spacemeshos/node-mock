package services

import (
	"math"
	"time"

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
	KnownPeers:    1,
	MinPeers:      1,
	MaxPeers:      1,
	IsSynced:      true,
	SyncedLayer:   1,
	CurrentLayer:  1,
	VerifiedLayer: 1,
}

var nodeError spacemesh.NodeError

var genesisTime = time.Now()

var currentEpoch uint64

var netID uint64 = 1
var layersPerEpoch uint64 = 10

var currentLayer = spacemesh.Layer{}

func updateSyncStatus() {
	if syncPosition >= len(syncStatusMock) {
		syncPosition = len(syncStatusMock) - 1
	}

	syncStatus = syncStatusMock[syncPosition]

	syncPosition++
}

func updateLayer() {
	currentLayer.Number++

	nodeStatus.CurrentLayer = currentLayer.Number

	if math.Mod(float64(nodeStatus.CurrentLayer), float64(layersPerEpoch)) == 0 {
		currentEpoch++
	}
}

// StatusLoadProducer -
func StatusLoadProducer() {
	for {
		updateSyncStatus()

		updateLayer()

		time.Sleep(10 * time.Second)
	}
}
