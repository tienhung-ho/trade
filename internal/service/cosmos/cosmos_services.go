package cosmosservice

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/tienhung-ho/mytoken/x/mytoken/types"
)

type CosmosInterface interface {
	CreateNewUser(alias, mnemonic, hdPath string) error
	MintAndSendTokens(owner string, denom string,
		amount int32, recipAddr string) *banktypes.MsgMintAndSendTokens
	BuildTx(msg sdk.Msg, denom, stake string, gasLimit uint64) (client.TxBuilder, error)
	SignTx(ctx context.Context, keyName string, txBuilder client.TxBuilder) error
	BroadcastTx(txBytes []byte) (sdk.TxResponse, error)
	GetAccount(address string) (sdk.AccountI, error)
	GenerateEntropy(bitSize int) ([]byte, error)
	GenerateMnemonic(entropy []byte) (string, error)
	GetAddressFromKeyName(keyName string) (string, error)
	EncodeTxBytes(txBuilder client.TxBuilder) ([]byte, error)
	UpdateTxFactoryWithLatestAccountInfo(address string) error
	UpdateTxFactoryAccountSequence(keyName string) error
	DeleteKey(alias string) error
	GetAccountSequence(address string) (uint64, error)
	SendTokens(owner string, denom string, amount int32, recipAddr string) *banktypes.MsgSendToken
	GetAllBalances(address string) (sdk.Coins, error)
	GetBalance(address, denom string) (*sdk.Coin, error)
	BroadcastSignedTx(signedTxBz []byte) (sdk.TxResponse, error)
	BuildUnsignedTx(
		msg sdk.Msg,
		denom, feeAmount string,
		gasLimit uint64,
		accountNumber uint64,
		sequence uint64,
		memo string,
	) ([]byte, error)
	GetTxFactoryAccAndSeq() (uint64, uint64, error)
}

type CosmosService struct {
	store CosmosInterface
}

func NewCosmosBiz(store CosmosInterface) *CosmosService {
	return &CosmosService{
		store: store,
	}
}

func (biz *CosmosService) BroadcastSignedTx(signedTxBase64 string) (sdk.TxResponse, error) {
	// 1. Decode base64
	signedTxBz, err := base64.StdEncoding.DecodeString(signedTxBase64)
	if err != nil {
		return sdk.TxResponse{}, fmt.Errorf("cannot decode base64 signed tx: %w", err)
	}

	// 2. Gọi method broadcast trong repo
	res, err := biz.store.BroadcastSignedTx(signedTxBz)
	if err != nil {
		return sdk.TxResponse{}, fmt.Errorf("cannot broadcast signed tx: %w", err)
	}

	return res, nil
}
