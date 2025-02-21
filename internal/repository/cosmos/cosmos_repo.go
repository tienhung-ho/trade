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

// BuildTx xây dựng transaction builder với message cần gửi.
func (ca *Cosmos) BuildTx(msg sdk.Msg) (client.TxBuilder, error) {
	txBuilder := ca.ClientCtx.TxConfig.NewTxBuilder()
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

func (ca *Cosmos) DeleteKey(alias string) error {
	return ca.Keyring.Delete(alias)
}
