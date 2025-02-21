package cosmosmodel

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc"
)

type CosmosConfig struct {
	ChainID        string
	NodeURI        string
	KeyringBackend string
	KeyringDir     string
	BroadcastMode  string
}

func NewCosmosConfig(chainID, nodeURI, keringBackend, keyringDir, broadcastMode string) *CosmosConfig {
	return &CosmosConfig{
		ChainID:        chainID,
		NodeURI:        nodeURI,
		KeyringBackend: keringBackend,
		KeyringDir:     keyringDir,
		BroadcastMode:  broadcastMode,
	}
}

type TxOptions struct {
	AccountNumber uint64    // Số thứ tự account trên chain
	Sequence      uint64    // Sequence của account
	Gas           uint64    // Lượng gas mặc định cho giao dịch
	GasAdjustment float64   // Hệ số điều chỉnh gas (ví dụ: 1.5)
	Fees          sdk.Coins // Phí giao dịch (có thể tạo bằng sdk.NewCoins(sdk.NewInt64Coin("stake", 200)))
	SignMode      string    // Chế độ ký, thường dùng mặc định của Cosmos SDK (ví dụ: "SIGN_MODE_DIRECT")
}

func NewTxConfig(accountNumber, sequence, gas uint64,
	gasAdjustment float64, fees sdk.Coins, signMode string) *TxOptions {
	return &TxOptions{
		AccountNumber: accountNumber,
		Sequence:      sequence,
		Gas:           gas,
		GasAdjustment: gasAdjustment,
		Fees:          fees,
		SignMode:      signMode,
	}

}

// AppContext chứa các thành phần cấu hình cần thiết cho các thao tác với Cosmos SDK:
// client context, tx factory, keyring, proto codec, …
type AppContext struct {
	ClientCtx  client.Context   // Client context để query, broadcast giao dịch, …
	TxFactory  tx.Factory       // Cấu hình factory cho việc xây dựng giao dịch
	Keyring    keyring.Keyring  // Quản lý key của các account
	ProtoCodec codec.ProtoCodec // Codec để mã hóa/giải mã dữ liệu
	GRPCConn   *grpc.ClientConn
}

func NewAppContext(clientCtx client.Context, txFactory tx.Factory,
	kering keyring.Keyring, protoCodec codec.ProtoCodec,

	conn *grpc.ClientConn) *AppContext {

	return &AppContext{
		ClientCtx:  clientCtx,
		TxFactory:  txFactory,
		Keyring:    kering,
		ProtoCodec: protoCodec,
		GRPCConn:   conn,
	}
}
