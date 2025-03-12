package cosmoshandler

import (
	"net/http"

	cosmosmodel "client/internal/model/cosmos"
	cosmosrepo "client/internal/repository/cosmos"
	cosmosservice "client/internal/service/cosmos"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/gin-gonic/gin"
)

// BroadCastTxRequest là struct để nhận body JSON
type BroadCastTxRequest struct {
	SignedTxBase64 string `json:"signed_tx"` // FE gửi "signed_tx" dưới dạng base64
}

// BroadCastTxSigned là handler cho route POST /tx/broadcast
func BroadCastTxSigned(appCtx *cosmosmodel.AppContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		// 1. Parse JSON body
		var req BroadCastTxRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 2. Tạo instance của service
		authQueryClient := authtypes.NewQueryClient(appCtx.GRPCConn)
		bankClient := banktypes.NewQueryClient(appCtx.GRPCConn)

		cosmosStore := cosmosrepo.NewCosmos(appCtx.ClientCtx, appCtx.TxFactory, appCtx.Keyring, authQueryClient, bankClient, &appCtx.ProtoCodec)
		cosmosBiz := cosmosservice.NewCosmosBiz(cosmosStore)
		// 3. Gọi service để broadcast
		res, err := cosmosBiz.BroadcastSignedTx(req.SignedTxBase64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 4. Trả về kết quả
		c.JSON(http.StatusOK, gin.H{
			"tx_hash":  res.TxHash,
			"raw_log":  res.RawLog,
			"height":   res.Height,
			"gas_used": res.GasUsed,
		})
	}
}
