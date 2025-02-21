package main

import (
	routerv1 "client/api/route/v1"
	cosmosmodel "client/internal/model/cosmos"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/joho/godotenv"
	banktypes "github.com/tienhung-ho/mytoken/x/mytoken/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// createNewUser tạo một tài khoản mới từ mnemonic và lưu vào keyring.
// Lưu ý: Trong môi trường production, bạn cần cẩn thận với việc quản lý mnemonic.
func createNewUser(kr keyring.Keyring, alias, mnemonic, hdPath string) error {
	// Tạo tài khoản mới với thuật toán Secp256k1
	_, err := kr.NewAccount(alias, mnemonic, "", hdPath, hd.Secp256k1)
	if err != nil {
		return fmt.Errorf("không thể tạo tài khoản %s: %w", alias, err)
	}
	log.Printf("Tài khoản mới đã được tạo thành công với alias '%s'.", alias)
	return nil
}

func mintAndSendTokens(owner sdk.AccountI, denom string, amount int32, recipAddr string) *banktypes.MsgMintAndSendTokens {
	return &banktypes.MsgMintAndSendTokens{
		Owner:     owner.GetAddress().String(),
		Denom:     denom,
		Amount:    amount,
		Recipient: recipAddr,
	}
}

// buildTx xây dựng transaction builder với message cần gửi.
func buildTx(clientCtx client.Context, txFactory tx.Factory, msg sdk.Msg) (client.TxBuilder, error) {
	txBuilder := clientCtx.TxConfig.NewTxBuilder()
	if err := txBuilder.SetMsgs(msg); err != nil {
		return nil, fmt.Errorf("lỗi khi set msg: %w", err)
	}
	// Ví dụ: set gas limit và fee nếu cần.
	txBuilder.SetGasLimit(300000)
	fees, err := sdk.ParseCoinsNormalized("500stake")
	if err != nil {
		return nil, fmt.Errorf("error parsing fees: %w", err)
	}
	txBuilder.SetFeeAmount(fees)
	return txBuilder, nil
}

// signTx ký giao dịch sử dụng key có alias keyName.
func signTx(ctx context.Context, txFactory tx.Factory, keyName string, txBuilder client.TxBuilder) error {
	if err := tx.Sign(ctx, txFactory, keyName, txBuilder, true); err != nil {
		return fmt.Errorf("lỗi khi ký giao dịch: %w", err)
	}
	return nil
}

// broadcastTx broadcast giao dịch và trả về kết quả.
func broadcastTx(clientCtx client.Context, txBytes []byte) (sdk.TxResponse, error) {
	res, err := clientCtx.BroadcastTxSync(txBytes)
	if err != nil {
		return sdk.TxResponse{}, fmt.Errorf("error broadcasting transaction: %w", err)
	}
	if res.Code != 0 {
		return *res, fmt.Errorf("transaction failed with code %d: %s", res.Code, res.RawLog)
	}
	return *res, nil
}

// getAccount truy vấn thông tin account từ auth module.
func getAccount(address string, authQueryClient authtypes.QueryClient, cdc codec.Codec) (sdk.AccountI, error) {
	req := &authtypes.QueryAccountRequest{
		Address: address,
	}
	resp, err := authQueryClient.Account(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("error getting account: %w", err)
	}
	var account sdk.AccountI
	err = cdc.InterfaceRegistry().UnpackAny(resp.Account, &account)
	if err != nil {
		return nil, fmt.Errorf("error unpacking account: %w", err)
	}
	return account, nil
}

func main() {
	// ----- Load .env và kết nối phụ trợ -----
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := mysqlConnection()
	if err != nil {
		log.Fatal("Error connecting to sql", err)
	}
	fmt.Println("Connected successfully to mysql", db)

	rdb := redisConnection()
	fmt.Println("Connected successfully doiredis", rdb)

	// Kết nối đến GRPC của node
	grpcConn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer grpcConn.Close()

	// Tạo cấu hình Cosmos (CosmosConfig)
	config := cosmosmodel.NewCosmosConfig(
		"mytoken",                    // ChainID
		"http://127.0.0.1:26657",     // NodeURI
		keyring.BackendTest,          // KeyringBackend
		"/home/tienhung-ho/.mytoken", // KeyringDir
		"block",                      // BroadcastMode
	)

	// Định nghĩa TxOptions (sử dụng giá trị tạm thời; sẽ cập nhật lại sau khi lấy được thông tin account)
	txOpts := cosmosmodel.NewTxConfig(
		1,      // AccountNumber (dummy)
		0,      // Sequence (dummy)
		200000, // Gas
		1.5,    // GasAdjustment
		sdk.NewCoins(sdk.NewInt64Coin("stake", 200)), // Fees
		"SIGN_MODE_DIRECT",                           // SignMode
	)

	// Khởi tạo AppContext
	appCtx, err := NewAppContext(config,
		txOpts, grpcConn)
	if err != nil {
		log.Fatalf("Không thể khởi tạo AppContext: %v", err)
	}

	// Truy vấn account của người gửi (ở đây sử dụng key "alice" đã được import từ chain)
	port := os.Getenv("PORT")
	r := routerv1.NewRouter(db, rdb, appCtx)
	if err := r.Run(port); err != nil {
	} // listen and serve (for windows "localhost:3000")
	return
}
