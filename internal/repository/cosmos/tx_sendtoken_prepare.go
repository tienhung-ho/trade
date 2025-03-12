package cosmosrepo

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BuildUnsignedTx xây dựng transaction CHƯA ký (offline signing).
// Trả về []byte (encode protobuf) để gửi cho client ký offline.
func (ca *Cosmos) BuildUnsignedTx(
	msg sdk.Msg,
	denom, feeAmount string,
	gasLimit uint64,
	accountNumber uint64,
	sequence uint64,
	memo string,
) ([]byte, error) {

	if msg == nil {
		return nil, fmt.Errorf("msg cannot be nil")
	}

	// Tạo TxBuilder mới từ ClientCtx.
	txBuilder := ca.ClientCtx.TxConfig.NewTxBuilder()

	// Set message
	if err := txBuilder.SetMsgs(msg); err != nil {
		return nil, fmt.Errorf("failed to set msg: %w", err)
	}

	// Set gas và fee
	txBuilder.SetGasLimit(gasLimit)

	feeStr := fmt.Sprintf("%s%s", feeAmount, denom)
	fees, err := sdk.ParseCoinsNormalized(feeStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing fee string %q: %w", feeStr, err)
	}
	txBuilder.SetFeeAmount(fees)

	// Set memo (tùy chọn)
	txBuilder.SetMemo(memo)

	// <--- Quan trọng --->
	// Không ký.
	// Thay vào đó, ta encode ra []byte để gửi cho client ký offline.
	// Kết quả: "unsignedTxBytes".

	unsignedTxBytes, err := ca.ClientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, fmt.Errorf("failed to encode unsigned tx: %w", err)
	}

	return unsignedTxBytes, nil
}

func (ca *Cosmos) BroadcastSignedTx(signedTxBz []byte) (sdk.TxResponse, error) {
	if len(signedTxBz) == 0 {
		return sdk.TxResponse{}, fmt.Errorf("signedTxBz is empty")
	}
	// Gửi transaction đã ký lên blockchain
	res, err := ca.ClientCtx.BroadcastTxSync(signedTxBz)
	if err != nil {
		return sdk.TxResponse{}, fmt.Errorf("error broadcasting transaction: %w", err)
	}
	if res.Code != 0 {
		return *res, fmt.Errorf("transaction failed with code %d: %s", res.Code, res.RawLog)
	}
	return *res, nil
}

func (ca *Cosmos) DecodeSignedTx(signedTxBz []byte) (sdk.Tx, error) {
	tx, err := ca.ClientCtx.TxConfig.TxDecoder()(signedTxBz)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signed tx: %w", err)
	}
	return tx, nil
}
