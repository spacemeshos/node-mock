package services

import (
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/spacemeshos/node-mock/spacemesh"
	"golang.org/x/exp/errors/fmt"
	"google.golang.org/grpc"
)

// GlobalStateService -
type GlobalStateService struct{}

func createAccount() (account spacemesh.Account) {
	key, _ := crypto.GenerateKey()

	account.Address.Address = crypto.PubkeyToAddress(key.PublicKey).Bytes()

	return
}

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

		time.Sleep(5 * time.Second)
	}
}

// RewardStream Rewards are computed by the protocol outside the STF but are a special case and are passed through the STF since they touch account balances.
func (s GlobalStateService) RewardStream(empty *empty.Empty, server spacemesh.GlobalStateService_RewardStreamServer) (err error) {
	return
}

// TransactionStateStream Transaction State - rejected pre-STF, or pending STF, or processed by STF
func (s GlobalStateService) TransactionStateStream(empty *empty.Empty, server spacemesh.GlobalStateService_TransactionStateStreamServer) (err error) {
	return
}

// TransactionReceiptStream Receipts - emitted after tx was processed by STF (or rejected before STF)
func (s GlobalStateService) TransactionReceiptStream(empty *empty.Empty, server spacemesh.GlobalStateService_TransactionReceiptStreamServer) (err error) {
	return
}

// InitGlobal -
func InitGlobal(s *grpc.Server) {
	spacemesh.RegisterGlobalStateServiceServer(s, GlobalStateService{})
}
