package services

import (
	"math/rand"
	"time"

	"github.com/spacemeshos/ed25519"

	"google.golang.org/grpc"

	"github.com/spacemeshos/go-spacemesh/common/types"

	"github.com/spacemeshos/node-mock/utils"

	v1 "github.com/spacemeshos/api/release/go/spacemesh/v1"
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
		Min          int
		Max          int
		MaxPerSecond uint64
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
	Smeshers struct {
		PoolSize int
	}
}

// Config -
var Config *Configuration

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

var currentLayerNumber = 0

var layerDuration = 500

var nodeError v1.NodeError

var genesisTime = time.Now()

var currentEpoch uint64

var currentLayer = &v1.Layer{}

type transactionInfo struct {
	Transaction v1.Transaction
	State       v1.TransactionState
	Receipt     v1.TransactionReceipt
}

// node
var syncStatusBus utils.Bus

//mesh
var layerBus utils.Bus

var layers []*v1.Layer
var activations []*v1.Activation
var blocks []*v1.Block
var transactions []*transactionInfo

// global
var rewardBus utils.Bus
var transactionStateBus utils.Bus
var transactionReceiptBus utils.Bus
var globalStateBus utils.Bus

var accounts []*v1.Account
var rewards []*v1.Reward

var smeshers []v1.SmesherId

func createAccount() (account v1.Account) {
	publicKey, _, _ := ed25519.GenerateKey(nil)

	//rand.New(rand.NewSource(time.Now().UnixNano()))

	accountID := v1.AccountId{
		Address: publicKey,
	}

	account = v1.Account{
		Address: &accountID,
	}

	accounts = append(accounts, &account)

	return
}

func getRandomBuffer(l int) []byte {
	result := make([]byte, l)

	rand.Read(result)

	return result
}

func getTransactionState() (result v1.TransactionState_TransactionState) {
	if rand.Float64() < 0.9 {
		result = v1.TransactionState_TRANSACTION_STATE_PROCESSED
	} else {
		result = (v1.TransactionState_TransactionState)(rand.Intn(4) + 1)
	}

	return
}

func transactionStateToResult(state v1.TransactionState_TransactionState) v1.TransactionReceipt_TransactionResult {
	switch state {
	case v1.TransactionState_TRANSACTION_STATE_UNSPECIFIED:
		return v1.TransactionReceipt_TRANSACTION_RESULT_UNSPECIFIED
	case v1.TransactionState_TRANSACTION_STATE_PROCESSED:
		return v1.TransactionReceipt_TRANSACTION_RESULT_EXECUTED
	case v1.TransactionState_TRANSACTION_STATE_INSUFFICIENT_FUNDS:
		return v1.TransactionReceipt_TRANSACTION_RESULT_INSUFFICIENT_FUNDS
	default:
		return v1.TransactionReceipt_TRANSACTION_RESULT_UNSPECIFIED
	}
}

func generateAtxID() (result types.AtxId) {
	for i := 0; i < len(result); i++ {
		result[i] = (byte)(rand.Intn(255))
	}

	return
}

func addrConvert(addr *v1.AccountId) (result types.Address) {
	for i := 0; i < len(result); i++ {
		result[i] = addr.Address[i]
	}

	return
}

/*func generateATX(layer uint64) []byte {
	account := getAccount()

	data := types.ActivationTxHeader{
		NIPSTChallenge: types.NIPSTChallenge{
			NodeId: types.NodeId{
				Key:          account.String(),
				VRFPublicKey: account.GetAddress(),
			},
			PubLayerIdx:          (types.LayerID)(layer),
			CommitmentMerkleRoot: getRandomBuffer(32),
		},
		Coinbase:      addrConvert(account),
		ActiveSetSize: (uint32)(rand.Intn(5)),
	}

	result := new(bytes.Buffer)

	encoder := json.NewEncoder(result)

	encoder.Encode(&data)

	return result.Bytes()
}*/

func generateTransaction(layerNumber uint64) (result *transactionInfo) {
	id := &v1.TransactionId{
		Id: getRandomBuffer(32),
	}

	txState := v1.TransactionState{
		Id:    id,
		State: getTransactionState(),
	}

	result = &transactionInfo{
		Transaction: v1.Transaction{
			Id: id,
		},
		State: v1.TransactionState{
			State: txState.State,
		},
		Receipt: v1.TransactionReceipt{
			Id:          txState.Id,
			Result:      transactionStateToResult(txState.State),
			GasUsed:     1000,
			LayerNumber: layerNumber,
		},
	}

	return
}

func generateTransactions(layer uint64) []*v1.Transaction {
	count := rand.Intn(Config.Transactions.Max-Config.Transactions.Min) + Config.Transactions.Min

	result := make([]*v1.Transaction, count)

	var txInfo *transactionInfo

	for i := 0; i < count; i++ {
		txInfo = generateTransaction(layer)

		transactionStateBus.Send(&txInfo.State)

		if txInfo.Receipt.Result != v1.TransactionReceipt_TRANSACTION_RESULT_UNSPECIFIED {
			transactionReceiptBus.Send(&txInfo.Receipt)
		}

		transactions = append(transactions, txInfo)

		result[i] = &txInfo.Transaction
	}

	return result
}

func generateBlocks(layer uint64) []*v1.Block {
	count := rand.Intn(Config.Blocks.Max-Config.Blocks.Min) + Config.Blocks.Min
	result := make([]*v1.Block, count)

	for i := 0; i < count; i++ {
		block := v1.Block{
			Id:           getRandomBuffer(32),
			Transactions: generateTransactions(layer),
		}

		blocks = append(blocks, &block)

		result[i] = &block
	}

	return result
}

func generateActivations(layer uint64) []*v1.Activation {
	count := rand.Intn(Config.Blocks.Max-Config.Blocks.Min) + Config.Blocks.Min
	result := make([]*v1.Activation, count)

	for i := 0; i < count; i++ {
		activation := v1.Activation{
			Id:    &v1.ActivationId{Id: getRandomBuffer(32)},
			Layer: layer,
		}

		activations = append(activations, &activation)

		result[i] = &activation
	}

	return result
}

func getAccount() *v1.AccountId {
	if len(accounts) < 100 {
		return createAccount().Address
	}
	return accounts[rand.Intn(len(accounts))].Address
}

func generateRewards(layer uint64) {
	count := rand.Intn(Config.Rewards.Max-Config.Rewards.Min) + Config.Rewards.Min

	for i := 0; i < count; i++ {
		reward := v1.Reward{
			Layer: layer,
			Total: &v1.Amount{
				Value: 10,
			},
			LayerReward: &v1.Amount{
				Value: 9,
			},
			LayerComputed: layer,
			Coinbase:      getAccount(),
			Smesher:       getSmesher(),
		}

		rewards = append(rewards, &reward)

		rewardBus.Send(&reward)
	}
}

func getLayerStatus(status v1.Layer_LayerStatus) int {
	var result int

	for _, v := range layers {
		if v.Status == status {
			result++
		}
	}

	return result
}

func updateLayers() {
	al := getLayerStatus(v1.Layer_LAYER_STATUS_APPROVED)

	if al > Config.Layers.MaxApproved {
		for _, v := range layers {
			if v.Status == v1.Layer_LAYER_STATUS_APPROVED {
				v.Status = v1.Layer_LAYER_STATUS_CONFIRMED

				layerBus.Send(v)

				generateRewards(v.Number)

				al--

				if al <= Config.Layers.MaxApproved {
					break
				}
			}
		}
	}

	for _, v := range layers {
		if v.Status == v1.Layer_LAYER_STATUS_UNSPECIFIED {
			v.Status = v1.Layer_LAYER_STATUS_APPROVED

			layerBus.Send(v)
		}
	}
}

func createLayer() {
	var layerNumber uint64

	if len(layers) != 0 {
		nodeStatus.TopLayer++

		layerNumber = nodeStatus.TopLayer
	} else {
		layerNumber = 0
	}

	currentLayer = &v1.Layer{
		Number:        layerNumber,
		Status:        v1.Layer_LAYER_STATUS_UNSPECIFIED,
		Hash:          getRandomBuffer(32),
		Blocks:        generateBlocks(layerNumber),
		Activations:   generateActivations(layerNumber),
		RootStateHash: getRandomBuffer(32),
	}

	if !nodeStatus.IsSynced {
		currentLayer.Status = v1.Layer_LAYER_STATUS_CONFIRMED
	}

	layerBus.Send(currentLayer)

	generateRewards(currentLayer.Number)

	layers = append(layers, currentLayer)
}

func startLoadProducer() {
	for {
		updateLayers()

		createLayer()

		if (!nodeStatus.IsSynced) && (len(layers) == Config.Threshold.Sync) {
			internalStatus = statusSynced

			nodeStatus.IsSynced = true

			syncStatusBus.Send(
				&v1.StatusStreamResponse{
					Status: &nodeStatus,
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

func createSmesher() *v1.SmesherId {
	var smesher v1.SmesherId

	accountID := getAccount()

	smesher.Id = accountID.Address

	return &smesher
}

func getSmesher() *v1.SmesherId {
	return &smeshers[rand.Intn(len(smeshers))]
}

func generateSmeshers() {
	for i := 0; i < Config.Smeshers.PoolSize; i++ {

		smeshers = append(smeshers, *createSmesher())
	}
}

// InitMocker -
func InitMocker(server *grpc.Server) {
	InitNode(server)
	InitMesh(server)
	InitGlobal(server)

	rand.Seed(time.Now().UnixNano())

	generateSmeshers()

	syncStatusBus.Init()
	layerBus.Init()
	globalStateBus.Init()

	rewardBus.Init()
	transactionStateBus.Init()
	transactionReceiptBus.Init()
}
