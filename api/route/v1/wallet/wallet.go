package walletv1

import (
	wallethandler "client/api/handler/wallet"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func WallterRouter(wallet *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	wallet.POST("/", wallethandler.CreateWalletHand(db, rdb))
}
