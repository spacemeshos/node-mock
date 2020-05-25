package services

import (
	"math/rand"
	"time"

	"github.com/spacemeshos/ed25519"

	"github.com/spacemeshos/node-mock/spacemesh"
	"github.com/spacemeshos/node-mock/utils"
	"google.golang.org/grpc"
)

// Configuration -
type Configuration struct {
	Version   string
	Build     string
	NetID     uint64
	RPCPort   uint
	Threshold struct {
		Sync   int
		Before int
		After  int
	}
	Transactions struct {
		Min int
		Max int
	}
	Blocks struct {
		Min int
		Max int
	}
	Rewards struct {
		Min int
		Max int
	}
	Layers struct {
		PerEpoch    uint64
		MaxApproved int
	}
}

// Config -
var Config Configuration

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

var currentLayer = &spacemesh.Layer{}

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
var rewardBus utils.Bus
var transactionStateBus utils.Bus
var transactionReceiptBus utils.Bus

var accounts []spacemesh.Account
var rewards []spacemesh.Reward

func createAccount() (account spacemesh.Account) {
	publicKey, _, _ := ed25519.GenerateKey(nil)

	//rand.New(rand.NewSource(time.Now().UnixNano()))

	accountID := spacemesh.AccountId{
		Address: publicKey,
	}

	account = spacemesh.Account{
		Address: &accountID,
	}

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
	count := rand.Intn(Config.Transactions.Max-Config.Transactions.Min) + Config.Transactions.Min

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
	count := rand.Intn(Config.Blocks.Max-Config.Blocks.Min) + Config.Blocks.Min
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

func getAccount() *spacemesh.AccountId {
	if len(accounts) < 100 {
		return createAccount().Address
	}
	return accounts[rand.Intn(len(accounts))].Address
}

func getSmesher() *spacemesh.SmesherId {
	var smesher spacemesh.SmesherId

	accountID := getAccount()

	smesher.Id = accountID.Address

	return &smesher
}

func generateRewards(layer uint64) {
	count := rand.Intn(Config.Rewards.Max-Config.Rewards.Min) + Config.Rewards.Min

	for i := 0; i < count; i++ {
		reward := spacemesh.Reward{
			Layer: layer,
			Total: &spacemesh.Amount{
				Value: 10,
			},
			LayerReward: &spacemesh.Amount{
				Value: 9,
			},
			LayerComputed: layer,
			Coinbase:      getAccount(),
			Smesher:       getSmesher(),
		}

		rewards = append(rewards, reward)

		rewardBus.Send(reward)
	}
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

	generateRewards(currentLayer.Number)

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

func updateLayers() {
	al := getLayerStatus(spacemesh.Layer_APPROVED)

	if al > Config.Layers.MaxApproved {
		for _, v := range layers {
			if v.Status == spacemesh.Layer_APPROVED {
				v.Status = spacemesh.Layer_CONFIRMED

				layerBus.Send(*v)

				generateRewards(v.Number)

				al--

				if al <= Config.Layers.MaxApproved {
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

		if (!nodeStatus.IsSynced) && (len(layers) == Config.Threshold.Sync) {
			syncStatusBus.Send(
				spacemesh.NodeSyncStatus{
					Status: spacemesh.NodeSyncStatus_SYNCED,
				},
			)

			nodeStatus.IsSynced = true
		}

		if len(layers) < Config.Threshold.Sync {
			time.Sleep((time.Duration)(Config.Threshold.Before) * time.Millisecond)
		} else {
			time.Sleep((time.Duration)(Config.Threshold.After) * time.Millisecond)
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
	rewardBus.Init()
	transactionStateBus.Init()
	transactionReceiptBus.Init()
}
