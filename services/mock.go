package services

import (
	"math/rand"
	"time"

	"google.golang.org/grpc"

//	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spacemeshos/go-spacemesh/log"

	"github.com/spacemeshos/node-mock/utils"
	"github.com/spacemeshos/node-mock/mock"
//	"github.com/spacemeshos/dash-backend/mock"

	"github.com/spacemeshos/api/release/go/spacemesh/v1"

)

const statusStopped = 0
const statusSyncing = 1
const statusSynced = 2

var internalStatus = statusStopped

var nodeStatus = v1.NodeStatus{
	ConnectedPeers: 0,
	IsSynced:       false,
	SyncedLayer:    0,
	TopLayer:       0,
	VerifiedLayer:  0,
}

var nodeError v1.NodeError

var network *mock.Network

// node
var syncStatusBus utils.Bus

//mesh
var layerBus utils.Bus

// global
var rewardBus utils.Bus
var transactionStateBus utils.Bus
var transactionReceiptBus utils.Bus
var globalStateBus utils.Bus

type proxy struct {
}

//    layerBus.Send(v)

func startLoadProducer() {
    for {
        network.Tick()
        time.Sleep(1 * time.Second)
    }
}

// InitMocker -
func InitMocker(server *grpc.Server, netId int, epochLayers int, maxTxs int, layerDuration int) {
    InitNode(server)
    InitMesh(server)
    InitGlobal(server)

    rand.Seed(time.Now().UnixNano())

    network = mock.NewNetwork(&proxy{}, netId, epochLayers, maxTxs, layerDuration)

    syncStatusBus.Init()
    layerBus.Init()
    globalStateBus.Init()

    rewardBus.Init()
    transactionStateBus.Init()
    transactionReceiptBus.Init()
}

func (proxy *proxy) SetNetworkInfo(netId uint64, genesisTime uint64, epochNumLayers uint64, maxTransactionsPerSecond uint64, layerDuration uint64) {
}

func (proxy *proxy) OnLayerChanged(layer *v1.Layer) {
    log.Info("proxy: OnLayerChanged(%v)", layer.Number)
    layerBus.Send(layer)
}

func (proxy *proxy) OnAccountChanged(account *v1.Account) {
    log.Info("proxy: OnAccountChanged")
    items := make([]*v1.GlobalStateDataItem, 1, 1)
    items[0] = &v1.GlobalStateDataItem{Data: &v1.GlobalStateDataItem_Account{Account: account}}
    globalStateBus.Send(&v1.GlobalStateStreamResponse{DataItem: items})
}

func (proxy *proxy) OnReward(reward *v1.Reward) {
    log.Info("proxy: OnReward")
    items := make([]*v1.GlobalStateDataItem, 1, 1)
    items[0] = &v1.GlobalStateDataItem{Data: &v1.GlobalStateDataItem_Reward{Reward: reward}}
    globalStateBus.Send(&v1.GlobalStateStreamResponse{DataItem: items})
}

func (proxy *proxy) OnTransactionReceipt(receipt *v1.TransactionReceipt) {
    log.Info("proxy: OnTransactionReceipt")
    items := make([]*v1.GlobalStateDataItem, 1, 1)
    items[0] = &v1.GlobalStateDataItem{Data: &v1.GlobalStateDataItem_Receipt{Receipt: receipt}}
    globalStateBus.Send(&v1.GlobalStateStreamResponse{DataItem: items})
}
