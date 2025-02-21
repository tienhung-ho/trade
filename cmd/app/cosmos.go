package main

import (
	cosmosmodel "client/internal/model/cosmos"
	"fmt"

	tmclient "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/tienhung-ho/mytoken/x/mytoken/types"
	"google.golang.org/grpc"
)

// NewAppContext khởi tạo và cấu hình các thành phần cần thiết dựa vào CosmosConfig và TxOptions.
func NewAppContext(cosmosConfig *cosmosmodel.CosmosConfig,
	txOpts *cosmosmodel.TxOptions, grpcConn *grpc.ClientConn) (*cosmosmodel.AppContext, error) {
	// Khởi tạo InterfaceRegistry và đăng ký các interface cần thiết
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	banktypes.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	cryptocodec.RegisterInterfaces(interfaceRegistry)

	// Tạo ProtoCodec dùng để mã hóa/giải mã
	protoCodec := codec.NewProtoCodec(interfaceRegistry)

	// Khởi tạo keyring sử dụng các thông số từ config
	kr, err := keyring.New(
		cosmosConfig.ChainID,
		cosmosConfig.KeyringBackend,
		cosmosConfig.KeyringDir,
		nil,
		protoCodec,
	)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi tạo keyring: %v", err)
	}

	// Khởi tạo client context với các thông số cấu hình
	clientCtx := client.Context{}.
		WithInterfaceRegistry(interfaceRegistry).
		WithKeyring(kr).
		WithNodeURI(cosmosConfig.NodeURI).
		WithCodec(protoCodec).
		WithChainID(cosmosConfig.ChainID).
		WithBroadcastMode(cosmosConfig.BroadcastMode).
		WithTxConfig(authtx.NewTxConfig(protoCodec, authtx.DefaultSignModes))

	// Tạo RPC client từ Tendermint/CometBFT và gán vào client context
	rpcClient, err := tmclient.New(cosmosConfig.NodeURI, "/websocket")
	if err != nil {
		return nil, fmt.Errorf("error creating RPC client: %v", err)
	}
	clientCtx = clientCtx.WithClient(rpcClient)

	// Khởi tạo txFactory với thông tin của account từ txOpts
	txFactory := tx.Factory{}.
		WithChainID(cosmosConfig.ChainID).
		WithKeybase(kr).
		WithTxConfig(clientCtx.TxConfig).
		WithAccountNumber(txOpts.AccountNumber).
		WithSequence(txOpts.Sequence).
		WithGas(txOpts.Gas).
		WithGasAdjustment(txOpts.GasAdjustment).
		WithFees(txOpts.Fees.String()).
		WithSignMode(authtx.DefaultSignModes[0]) // Nếu muốn chuyển sang chế độ ký khác thì mapping từ txOpts.SignMode

	// Tạo AppContext sử dụng hàm NewAppContext đã khai báo trong package cosmosmodel
	appCtx := cosmosmodel.NewAppContext(clientCtx, txFactory, kr, *protoCodec, grpcConn)

	return appCtx, nil
}
