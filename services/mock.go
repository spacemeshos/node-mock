package services

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"github.com/spacemeshos/ed25519"

	"github.com/spacemeshos/node-mock/spacemesh"
	"github.com/spacemeshos/node-mock/utils"
	"google.golang.org/grpc"
)

const mockVersion = "0.0.1"
const mockBuild = "1"

var nodeStatus = spacemesh.NodeStatus{
	KnownPeers:    1,
	MinPeers:      1,
	MaxPeers:      1,
	IsSynced:      false,
	SyncedLayer:   0,
	CurrentLayer:  0,
	VerifiedLayer: 0,
}

var nodeError spacemesh.NodeError

var genesisTime = time.Now()

var currentEpoch uint64

var netID uint64 = 1
var layersPerEpoch uint64 = 10

var minBlocks = 1
var maxBlocks = 10

var minTransactions = 1
var maxTransactions = 10

var producerIntervalBS = 1  // seconds
var producerIntervalAS = 15 // seconds

var maxApprovedLayers = 3

var currentLayer = &spacemesh.Layer{}

var syncThreshold = 10

type transactionInfo struct {
	Transaction spacemesh.Transaction
	State       spacemesh.TransactionState
	Receipt     spacemesh.TransactionReceipt
}

// node
var syncStatusBus utils.Bus

//mesh
var layerBus utils.Bus

var layers []*spacemesh.Layer
var blocks []spacemesh.Block
var transactions []transactionInfo

// global
var accounts []spacemesh.Account
var rewards []spacemesh.Reward
var transactionStateBus utils.Bus
var transactionReceiptBus utils.Bus

func createAccount() (account spacemesh.Account) {
	publicKey, _, _ := ed25519.GenerateKey(nil)

	account.Address.Address = publicKey

	accounts = append(accounts, account)

	return
}

func getRandomBuffer(l int) []byte {
	result := make([]byte, l)

	rand.Read(result)

	return result
}

func getTransactionState() (result spacemesh.TransactionState_TransactionStateType) {
	if rand.Float64() < 0.9 {
		result = spacemesh.TransactionState_PROCESSED
	} else {
		result = (spacemesh.TransactionState_TransactionStateType)(rand.Intn(3) + 2)
	}

	return
}

func transactionStateToResult(state spacemesh.TransactionState_TransactionStateType) spacemesh.TransactionReceipt_TransactionResult {
	switch state {
	case spacemesh.TransactionState_UNDEFINED:
		return spacemesh.TransactionReceipt_UNDEFINED
	case spacemesh.TransactionState_UNKNOWN:
		return spacemesh.TransactionReceipt_UNKNOWN
	case spacemesh.TransactionState_PROCESSED:
		return spacemesh.TransactionReceipt_EXECUTED
	case spacemesh.TransactionState_INSUFFICIENT_FUNDS:
		return spacemesh.TransactionReceipt_INSUFFICIENT_FUNDS
	default:
		return spacemesh.TransactionReceipt_UNDEFINED
	}
}

func generateTransactions(layer uint64) []*spacemesh.Transaction {
	count := rand.Intn(maxTransactions-minTransactions) + minTransactions

	result := make([]*spacemesh.Transaction, count)

	for i := 0; i < count; i++ {
		tx := spacemesh.Transaction{
			Type: spacemesh.Transaction_SIMPLE,
			Id: &spacemesh.TransactionId{
				Id: getRandomBuffer(32),
			},
		}

		txState := getTransactionState()

		txInfo := transactionInfo{
			Transaction: tx,
			State: spacemesh.TransactionState{
				State: txState,
			},
			Receipt: spacemesh.TransactionReceipt{
				Id:          tx.Id,
				Result:      transactionStateToResult(txState),
				GasUsed:     1000,
				LayerNumber: layer,
			},
		}

		transactionStateBus.Send(txInfo.State)

		if txInfo.Receipt.Result != spacemesh.TransactionReceipt_UNKNOWN {
			transactionReceiptBus.Send(txInfo.Receipt)
		}

		transactions = append(
			transactions,
			txInfo,
		)

		result[i] = &tx
	}

	return result
}

func generateBlocks(layer uint64) []*spacemesh.Block {
	count := rand.Intn(maxBlocks-minBlocks) + minBlocks

	result := make([]*spacemesh.Block, count)

	for i := 0; i < count; i++ {
		block := spacemesh.Block{
			Id:           getRandomBuffer(32),
			Transactions: generateTransactions(layer),
		}

		blocks = append(blocks, block)

		result[i] = &block
	}

	return result
}

func createLayer() {
	var layerNumber uint64

	if len(layers) != 0 {
		nodeStatus.CurrentLayer++

		layerNumber = nodeStatus.CurrentLayer
	} else {
		layerNumber = 0
	}

	currentLayer = &spacemesh.Layer{
		Number:        layerNumber,
		Status:        spacemesh.Layer_UNKNOWN,
		Hash:          getRandomBuffer(32),
		Blocks:        generateBlocks(layerNumber),
		RootStateHash: getRandomBuffer(32),
	}

	if !nodeStatus.IsSynced {
		currentLayer.Status = spacemesh.Layer_CONFIRMED
	}

	syncStatusBus.Send(
		spacemesh.NodeSyncStatus{
			Status: spacemesh.NodeSyncStatus_NEW_TOP_LAYER,
		},
	)

	layerBus.Send(*currentLayer)

	layers = append(layers, currentLayer)
}

func getLayerStatus(status spacemesh.Layer_LayerStatus) int {
	var result int

	for _, v := range layers {
		if v.Status == status {
			result++
		}
	}

	return result
}

func printLayers() {
	fmt.Println("-- Layers --")

	for _, v := range layers {
		fmt.Printf("%d - %s - %s\n", v.GetNumber(), hex.EncodeToString(v.GetHash()), v.GetStatus().String())
	}

	fmt.Println("-- ------ --")
}

func updateLayers() {
	al := getLayerStatus(spacemesh.Layer_APPROVED)

	if al > maxApprovedLayers {
		for _, v := range layers {
			if v.Status == spacemesh.Layer_APPROVED {
				v.Status = spacemesh.Layer_CONFIRMED

				layerBus.Send(*v)

				al--

				if al <= maxApprovedLayers {
					break
				}
			}
		}
	}

	for _, v := range layers {
		if v.Status == spacemesh.Layer_UNKNOWN {
			v.Status = spacemesh.Layer_APPROVED

			layerBus.Send(*v)
		}
	}
}

func startLoadProducer() {
	for {
		updateLayers()

		createLayer()

		if (!nodeStatus.IsSynced) && (len(layers) == syncThreshold) {
			syncStatusBus.Send(
				spacemesh.NodeSyncStatus{
					Status: spacemesh.NodeSyncStatus_SYNCED,
				},
			)

			nodeStatus.IsSynced = true
		}

		if len(layers) < syncThreshold {
			time.Sleep((time.Duration)(producerIntervalBS) * time.Second)
		} else {
			time.Sleep((time.Duration)(producerIntervalAS) * time.Second)
		}
	}
}

// InitMocker -
func InitMocker(server *grpc.Server) {
	InitNode(server)
	InitMesh(server)
	InitGlobal(server)

	rand.Seed(time.Now().UnixNano())

	syncStatusBus.Init()
	layerBus.Init()
	transactionStateBus.Init()
	transactionReceiptBus.Init()
}
