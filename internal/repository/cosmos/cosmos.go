package cosmosrepo

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// CosmosAuth chứa các dependency cần thiết cho các thao tác với Cosmos SDK.
type Cosmos struct {
	ClientCtx       client.Context
	TxFactory       tx.Factory
	Keyring         keyring.Keyring
	AuthQueryClient authtypes.QueryClient
	BankQueryClient banktypes.QueryClient
	Codec           codec.Codec
}

// NewCosmosAuth khởi tạo một instance của CosmosAuth.
func NewCosmos(
	clientCtx client.Context,
	txFactory tx.Factory,
	kr keyring.Keyring,
	authQueryClient authtypes.QueryClient,
	bankQueryClient banktypes.QueryClient,
	cdc codec.Codec) *Cosmos {
	return &Cosmos{
		ClientCtx:       clientCtx,
		TxFactory:       txFactory,
		Keyring:         kr,
		AuthQueryClient: authQueryClient,
		BankQueryClient: bankQueryClient,
		Codec:           cdc,
	}
}
