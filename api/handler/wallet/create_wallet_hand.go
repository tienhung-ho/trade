package wallethandler

import (
	"client/internal/common/apperrors"
	"client/internal/common/appresponses"
	walletmodel "client/internal/model/mysql/wallet"
	walletrepo "client/internal/repository/mysql/wallets"
	walletbus "client/internal/service/wallet"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func CreateWalletHand(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		var data walletmodel.UserWallet
		if err := c.ShouldBindJSON(&data); err != nil {

			log.Print(err)
			c.JSON(http.StatusBadRequest, apperrors.ErrInternal(err))
			c.Abort()
			return
		}

		store := walletrepo.NewMySQLWallet(db)
		biz := walletbus.NewWalletBiz(store, db)

		recordID, err := biz.CreateWallet(&data)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, appresponses.NewDataResponse(recordID, "create wallet successfully"))
	}
}
