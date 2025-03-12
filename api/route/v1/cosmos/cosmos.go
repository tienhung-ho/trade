package cosmosv1

import (
	cosmoshandler "client/api/handler/cosmos"
	cosmosmodel "client/internal/model/cosmos"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

func CosmosRouter(cosmos *gin.RouterGroup, db *gorm.DB, rdb *redis.Client, appCtx *cosmosmodel.AppContext) {
	cosmos.POST("/broadcast", cosmoshandler.BroadCastTxSigned(appCtx))
	cosmos.POST("/faucet", cosmoshandler.FaucetHandler(appCtx))
}
