package cosmosservice

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/tienhung-ho/mytoken/x/mytoken/types"
)

type CosmosInterface interface {
	CreateNewUser(alias, mnemonic, hdPath string) error
	MintAndSendTokens(owner string, denom string,
		amount int32, recipAddr string) *banktypes.MsgMintAndSendTokens
	BuildTx(msg sdk.Msg) (client.TxBuilder, error)
	SignTx(ctx context.Context, keyName string, txBuilder client.TxBuilder) error
	BroadcastTx(txBytes []byte) (sdk.TxResponse, error)
	GetAccount(address string) (sdk.AccountI, error)
	GenerateEntropy(bitSize int) ([]byte, error)
	GenerateMnemonic(entropy []byte) (string, error)
	GetAddressFromKeyName(keyName string) (string, error)
	EncodeTxBytes(txBuilder client.TxBuilder) ([]byte, error)
	UpdateTxFactoryAccountSequence(keyName string) error
	DeleteKey(alias string) error
}

type CosmosService struct {
	store CosmosInterface
}

func NewCosmosBiz(store CosmosInterface) *CosmosService {
	return &CosmosService{
		store: store,
	}
}
