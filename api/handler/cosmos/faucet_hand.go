package cosmoshandler

import (
	cosmosmodel "client/internal/model/cosmos"
	cosmosrepo "client/internal/repository/cosmos"
	cosmosservice "client/internal/service/cosmos"
	"net/http"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gin-gonic/gin"
)

type FaucetRequest struct {
	Address string `json:"address"` // Địa chỉ user mới tạo
}

type FaucetResponse struct {
	TxHash string `json:"tx_hash"`
	// Hoặc thêm code, log, ...
}

func FaucetHandler(appCtx *cosmosmodel.AppContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req FaucetRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Gọi hàm biz
		// 2. Tạo instance của service
		authQueryClient := authtypes.NewQueryClient(appCtx.GRPCConn)
		bankClient := banktypes.NewQueryClient(appCtx.GRPCConn)

		cosmosStore := cosmosrepo.NewCosmos(appCtx.ClientCtx, appCtx.TxFactory, appCtx.Keyring, authQueryClient, bankClient, &appCtx.ProtoCodec)
		cosmosBiz := cosmosservice.NewFaucetBiz(cosmosStore)
		txHash, err := cosmosBiz.FaucetToken(c.Request.Context(), req.Address)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, FaucetResponse{TxHash: txHash})

	}
}
