package services

import (
	"google.golang.org/grpc"

	context "context"

	v1 "github.com/spacemeshos/api/release/go/spacemesh/v1"
)

// GlobalStateService -
type GlobalStateService struct{}

// GlobalStateHash - Latest computed global state - layer and its root hash
func (s GlobalStateService) GlobalStateHash(ctx context.Context, request *v1.GlobalStateHashRequest) (*v1.GlobalStateHashResponse, error) {
	return &v1.GlobalStateHashResponse{}, nil
}

// Account info in the current global state.
func (s GlobalStateService) Account(ctx context.Context, request *v1.AccountRequest) (*v1.AccountResponse, error) {
	return &v1.AccountResponse{}, nil
}

// AccountDataQuery - Query for account related data such as rewards, tx receipts and account info
//
// Note: it might be too expensive to add a param for layer to get these results from
// as it may require indexing all global state changes per account by layer.
// If it is possible to index by layer then we should add param start_layer to
// AccountDataParams. Currently it will return data from genesis.
func (s GlobalStateService) AccountDataQuery(ctx context.Context, request *v1.AccountDataQueryRequest) (*v1.AccountDataQueryResponse, error) {
	return &v1.AccountDataQueryResponse{}, nil
}

// SmesherDataQuery Query for smesher data. Currently returns smesher rewards.
// Note: Not supporting start_layer yet as it may require to index all rewards by
// smesher and by layer id or allow for queries from a layer and later....
func (s GlobalStateService) SmesherDataQuery(ctx context.Context, request *v1.SmesherDataQueryRequest) (*v1.SmesherDataQueryResponse, error) {
	return &v1.SmesherDataQueryResponse{}, nil
}

// AccountStream Account changes (e.g., balance and counter/nonce changes).
/*func (s GlobalStateService) AccountStream(emty *empty.Empty, server v1.GlobalStateService_AccountStreamServer) (err error) {
	for {
		account := createAccount()

		err = server.Send(&account)
		if err != nil {
			fmt.Printf("AccountStream(ERROR): %v\n", err)

			return
		}

		fmt.Printf("AccountStream(OK): %v\n", account)

		time.Sleep(15 * time.Second)
	}
}*/

// AccountDataStream Get a stream of account related changes such as account balance change,
// tx receipts and rewards
func (s GlobalStateService) AccountDataStream(request *v1.AccountDataStreamRequest, server v1.GlobalStateService_AccountDataStreamServer) error {
	return nil
}

// RewardStream Rewards are computed by the protocol outside the STF but are a special case and are passed through the STF since they touch account balances.
/*func (s GlobalStateService) RewardStream(empty *empty.Empty, server v1.GlobalStateService_RewardStreamServer) (err error) {
	rewardChan, cookie := rewardBus.Register()
	defer rewardBus.Delete(cookie)

	for {
		select {
		case msg := <-rewardChan:
			reward := msg.(v1.Reward)

			err = server.Send(&reward)
			if err != nil {
				fmt.Printf("RewardStream(ERROR): %v\n", err)

				return
			}

			fmt.Printf("RewardStream(OK): %d - %d - %s\n", reward.GetLayer(), reward.GetLayerComputed(), reward.GetTotal().String())
		}
	}
}*/

// SmesherRewardStream Rewards awarded to a smesher id
func (s GlobalStateService) SmesherRewardStream(request *v1.SmesherRewardStreamRequest, server v1.GlobalStateService_SmesherRewardStreamServer) error {
	return nil
}

// AppEventStream App Events - emitted by app methods impl code trigged by an
// app transaction
func (s GlobalStateService) AppEventStream(request *v1.AppEventStreamRequest, server v1.GlobalStateService_AppEventStreamServer) error {
	return nil
}

// GlobalStateStream New global state computed for a layer by the STF
func (s GlobalStateService) GlobalStateStream(request *v1.GlobalStateStreamRequest, server v1.GlobalStateService_GlobalStateStreamServer) error {
	return nil
}

// TransactionStateStream Transaction State - rejected pre-STF, or pending STF, or processed by STF
/*func (s GlobalStateService) TransactionStateStream(empty *empty.Empty, server v1.GlobalStateService_TransactionStateStreamServer) (err error) {
	txStateChan, cookie := transactionStateBus.Register()
	defer transactionStateBus.Delete(cookie)

	for {
		select {
		case msg := <-txStateChan:
			txState := msg.(v1.TransactionState)

			err = server.Send(&txState)
			if err != nil {
				fmt.Printf("TransactionStateStream(ERROR): %v\n", err)

				return
			}

			fmt.Printf("TransactionStateStream(OK): %s\n", txState.String())
		}
	}
}*/

// TransactionReceiptStream Receipts - emitted after tx was processed by STF (or rejected before STF)
/*func (s GlobalStateService) TransactionReceiptStream(empty *empty.Empty, server v1.GlobalStateService_TransactionReceiptStreamServer) (err error) {
	txReceiptChan, cookie := transactionReceiptBus.Register()
	defer transactionReceiptBus.Delete(cookie)

	for {
		select {
		case msg := <-txReceiptChan:
			txReceipt := msg.(v1.TransactionReceipt)

			err = server.Send(&txReceipt)
			if err != nil {
				fmt.Printf("TransactionReceiptStream(ERROR): %v\n", err)

				return
			}

			fmt.Printf("TransactionReceiptStream(OK): %d - %s - %s\n",
				txReceipt.GetLayerNumber(),
				hex.EncodeToString(txReceipt.Id.Id),
				txReceipt.GetResult().String(),
			)
		}
	}
}*/

// InitGlobal -
func InitGlobal(s *grpc.Server) {
	v1.RegisterGlobalStateServiceServer(s, GlobalStateService{})
}
