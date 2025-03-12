package cosmosrepo

import (
	"context"
	"fmt"
	"log"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	csBanktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	banktypes "github.com/tienhung-ho/mytoken/x/mytoken/types"
)

// CreateNewUser tạo một tài khoản mới từ mnemonic và lưu vào keyring.
// Lưu ý: Trong môi trường production, cần cẩn thận khi quản lý mnemonic.
func (ca *Cosmos) CreateNewUser(alias, mnemonic, hdPath string) error {

	_, err := ca.Keyring.NewAccount(alias, mnemonic, "", hdPath, hd.Secp256k1)
	if err != nil {
		return fmt.Errorf("không thể tạo tài khoản %s: %w", alias, err)
	}
	log.Printf("Tài khoản mới đã được tạo thành công với alias '%s'.", alias)
	return nil
}

// MintAndSendTokens tạo message mint và gửi token.
func (ca *Cosmos) MintAndSendTokens(owner string, denom string, amount int32, recipAddr string) *banktypes.MsgMintAndSendTokens {
	return &banktypes.MsgMintAndSendTokens{
		Owner:     owner,
		Denom:     denom,
		Amount:    amount,
		Recipient: recipAddr,
	}
}

func (ca *Cosmos) GetAccountSequence(address string) (uint64, error) {
	// Query account info from blockchain
	accountInfo, err := ca.ClientCtx.AccountRetriever.GetAccount(ca.ClientCtx, sdk.MustAccAddressFromBech32(address))
	if err != nil {
		return 0, fmt.Errorf("failed to get account info: %w", err)
	}
	return accountInfo.GetSequence(), nil
}

func (ca *Cosmos) GetAllBalances(address string) (sdk.Coins, error) {
	req := &csBanktypes.QueryAllBalancesRequest{Address: address}
	resp, err := ca.BankQueryClient.AllBalances(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("error querying all balances: %w", err)
	}
	return resp.Balances, nil
}

func (ca *Cosmos) GetBalance(address, denom string) (*sdk.Coin, error) {
	req := &csBanktypes.QueryBalanceRequest{
		Address: address,
		Denom:   denom,
	}
	resp, err := ca.BankQueryClient.Balance(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("error querying balance for denom %s: %w", denom, err)
	}
	return resp.Balance, nil
}

// MintAndSendTokens tạo message mint và gửi token.
func (ca *Cosmos) SendTokens(owner string, denom string, amount int32, recipAddr string) *banktypes.MsgSendToken {
	return &banktypes.MsgSendToken{
		Owner:     owner,
		Denom:     denom,
		Amount:    amount,
		Recipient: recipAddr,
	}
}

// BuildTx xây dựng transaction builder với message cần gửi.
// BuildTx constructs a transaction builder with the provided message, fee and gas limit.
// The fee is built by concatenating the stake amount (as a string) with the denom,
// e.g., if stake = "500" and denom = "citcoin", then feeStr becomes "500citcoin".
func (ca *Cosmos) BuildTx(msg sdk.Msg, denom, stake string, gasLimit uint64) (client.TxBuilder, error) {
	if msg == nil {
		return nil, fmt.Errorf("msg cannot be nil")
	}

	// Create a new TxBuilder from ClientCtx.
	txBuilder := ca.ClientCtx.TxConfig.NewTxBuilder()
	if err := txBuilder.SetMsgs(msg); err != nil {
		return nil, fmt.Errorf("failed to set msg: %w", err)
	}

	// Set the provided gas limit.
	txBuilder.SetGasLimit(gasLimit)

	// Build fee string from stake and denom, e.g., "500citcoin".
	feeStr := fmt.Sprintf("%s%s", stake, denom)
	fees, err := sdk.ParseCoinsNormalized(feeStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing fee string %s: %w", feeStr, err)
	}
	txBuilder.SetFeeAmount(fees)

	return txBuilder, nil
}

// SignTx ký giao dịch sử dụng key có alias keyName.
func (ca *Cosmos) SignTx(ctx context.Context, keyName string, txBuilder client.TxBuilder) error {
	if err := tx.Sign(ctx, ca.TxFactory, keyName, txBuilder, true); err != nil {
		return fmt.Errorf("lỗi khi ký giao dịch: %w", err)
	}
	return nil
}

// BroadcastTx broadcast giao dịch và trả về kết quả.
func (ca *Cosmos) BroadcastTx(txBytes []byte) (sdk.TxResponse, error) {
	res, err := ca.ClientCtx.BroadcastTxSync(txBytes)
	if err != nil {
		return sdk.TxResponse{}, fmt.Errorf("error broadcasting transaction: %w", err)
	}
	if res.Code != 0 {
		return *res, fmt.Errorf("transaction failed with code %d: %s", res.Code, res.RawLog)
	}
	return *res, nil
}

// GetAccount truy vấn thông tin account từ auth module.
func (ca *Cosmos) GetAccount(address string) (sdk.AccountI, error) {
	req := &authtypes.QueryAccountRequest{
		Address: address,
	}
	resp, err := ca.AuthQueryClient.Account(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("error getting account: %w", err)
	}
	var account sdk.AccountI
	err = ca.Codec.InterfaceRegistry().UnpackAny(resp.Account, &account)
	if err != nil {
		return nil, fmt.Errorf("error unpacking account: %w", err)
	}
	return account, nil
}

func (ca *Cosmos) GetAddressFromKeyName(keyName string) (string, error) {
	// Lấy thông tin key (KeyInfo) từ keyring dựa vào alias (keyName)
	info, err := ca.Keyring.Key(keyName)
	if err != nil {
		return "", fmt.Errorf("cannot find key with alias '%s': %w", keyName, err)
	}

	// Lấy ra địa chỉ (Address) từ KeyInfo
	addr, err := info.GetAddress()
	if err != nil {
		return "", fmt.Errorf("cannot get address from key '%s': %w", keyName, err)
	}

	return addr.String(), nil
}

func (ca *Cosmos) EncodeTxBytes(txBuilder client.TxBuilder) ([]byte, error) {
	// Lấy đối tượng transaction (kiểu proto hoặc signing.Tx)
	tx := txBuilder.GetTx()

	// Dùng TxEncoder trong ClientCtx để mã hoá thành []byte
	txBytes, err := ca.ClientCtx.TxConfig.TxEncoder()(tx)
	if err != nil {
		return nil, fmt.Errorf("cannot encode transaction: %w", err)
	}
	return txBytes, nil
}

// Lấy account từ alias keyName, rồi query on-chain để biết accountNumber, sequence
func (ca *Cosmos) UpdateTxFactoryAccountSequence(keyName string) error {
	// 1) Lấy keyInfo từ keyring
	info, err := ca.Keyring.Key(keyName)
	if err != nil {
		return fmt.Errorf("cannot find key '%s': %w", keyName, err)
	}
	addr, err := info.GetAddress()
	if err != nil {
		return fmt.Errorf("cannot get address from key '%s': %w", keyName, err)
	}

	// 2) Query on-chain account
	acc, err := ca.GetAccount(addr.String())
	if err != nil {
		return fmt.Errorf("cannot get account from chain: %w", err)
	}

	// 3) Lấy accountNumber, sequence
	accNumber := acc.GetAccountNumber()
	seq := acc.GetSequence()

	// 4) Cập nhật TxFactory
	ca.TxFactory = ca.TxFactory.
		WithAccountNumber(accNumber).
		WithSequence(seq)

	return nil
}

// Thêm hàm này vào struct Cosmos của bạn
func (ca *Cosmos) GetLatestSequence(address string) (uint64, error) {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return 0, fmt.Errorf("địa chỉ không hợp lệ: %w", err)
	}

	// Sử dụng AccountRetriever để lấy thông tin tài khoản trực tiếp từ blockchain
	acc, err := ca.ClientCtx.AccountRetriever.GetAccount(ca.ClientCtx, addr)
	if err != nil {
		return 0, fmt.Errorf("không thể lấy thông tin tài khoản: %w", err)
	}

	return acc.GetSequence(), nil
}

// Cập nhật TxFactory với sequence mới nhất
// UpdateTxFactoryWithLatestAccountInfo lấy account number và sequence mới nhất từ blockchain
// sau đó cập nhật vào TxFactory. Hàm này đảm bảo rằng cả account number và sequence đều được cập nhật.
func (ca *Cosmos) UpdateTxFactoryWithLatestAccountInfo(address string) error {
	// Chuyển đổi địa chỉ từ bech32 sang sdk.AccAddress
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}

	// Lấy thông tin account trực tiếp từ blockchain thông qua AccountRetriever
	acc, err := ca.ClientCtx.AccountRetriever.GetAccount(ca.ClientCtx, addr)
	if err != nil {
		return fmt.Errorf("failed to get account info: %w", err)
	}

	// Logging để kiểm tra thông tin account
	log.Printf("Updating TxFactory: accountNumber=%d, sequence=%d", acc.GetAccountNumber(), acc.GetSequence())

	// Cập nhật TxFactory với account number và sequence mới nhất
	ca.TxFactory = ca.TxFactory.
		WithAccountNumber(acc.GetAccountNumber()).
		WithSequence(acc.GetSequence())

	return nil
}

func (ca *Cosmos) DeleteKey(alias string) error {
	return ca.Keyring.Delete(alias)
}

func (ca *Cosmos) GetTxFactoryAccAndSeq() (uint64, uint64, error) {
	// có thể kiểm tra nil, ...
	accNum := ca.TxFactory.AccountNumber()
	seq := ca.TxFactory.Sequence()
	return accNum, seq, nil
}
