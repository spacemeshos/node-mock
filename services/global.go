package services

import (
	"encoding/hex"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/spacemeshos/node-mock/spacemesh"
	"golang.org/x/exp/errors/fmt"
	"google.golang.org/grpc"
)

// GlobalStateService -
type GlobalStateService struct{}

// AccountStream Account changes (e.g., balance and counter/nonce changes).
func (s GlobalStateService) AccountStream(emty *empty.Empty, server spacemesh.GlobalStateService_AccountStreamServer) (err error) {
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
}

// RewardStream Rewards are computed by the protocol outside the STF but are a special case and are passed through the STF since they touch account balances.
func (s GlobalStateService) RewardStream(empty *empty.Empty, server spacemesh.GlobalStateService_RewardStreamServer) (err error) {
	rewardChan, cookie := rewardBus.Register()
	defer rewardBus.Delete(cookie)

	for {
		select {
		case msg := <-rewardChan:
			reward := msg.(spacemesh.Reward)

			err = server.Send(&reward)
			if err != nil {
				fmt.Printf("RewardStream(ERROR): %v\n", err)

				return
			}

			fmt.Printf("RewardStream(OK): %d - %d - %s\n", reward.GetLayer(), reward.GetLayerComputed(), reward.GetTotal().String())
		}
	}
}

// TransactionStateStream Transaction State - rejected pre-STF, or pending STF, or processed by STF
func (s GlobalStateService) TransactionStateStream(empty *empty.Empty, server spacemesh.GlobalStateService_TransactionStateStreamServer) (err error) {
	txStateChan, cookie := transactionStateBus.Register()
	defer transactionStateBus.Delete(cookie)

	for {
		select {
		case msg := <-txStateChan:
			txState := msg.(spacemesh.TransactionState)

			err = server.Send(&txState)
			if err != nil {
				fmt.Printf("TransactionStateStream(ERROR): %v\n", err)

				return
			}

			fmt.Printf("TransactionStateStream(OK): %s\n", txState.String())
		}
	}
}

// TransactionReceiptStream Receipts - emitted after tx was processed by STF (or rejected before STF)
func (s GlobalStateService) TransactionReceiptStream(empty *empty.Empty, server spacemesh.GlobalStateService_TransactionReceiptStreamServer) (err error) {
	txReceiptChan, cookie := transactionReceiptBus.Register()
	defer transactionReceiptBus.Delete(cookie)

	for {
		select {
		case msg := <-txReceiptChan:
			txReceipt := msg.(spacemesh.TransactionReceipt)

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
}

// InitGlobal -
func InitGlobal(s *grpc.Server) {
	spacemesh.RegisterGlobalStateServiceServer(s, GlobalStateService{})
}
