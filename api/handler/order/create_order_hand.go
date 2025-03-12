package orderhandler

import (
	"client/internal/common/apperrors"
	"client/internal/common/appresponses"
	cosmosmodel "client/internal/model/cosmos"
	ordermodel "client/internal/model/mysql/order"
	cosmosrepo "client/internal/repository/cosmos"
	orderrepo "client/internal/repository/mysql/order"
	orderitemrepo "client/internal/repository/mysql/order_item"
	productrepo "client/internal/repository/mysql/product"
	userrepo "client/internal/repository/mysql/user"
	orderbus "client/internal/service/order"
	"log"
	"net/http"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func CreateOrderhandler(db *gorm.DB, rdb *redis.Client, appCtx *cosmosmodel.AppContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var data ordermodel.CreateOrder

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusInternalServerError, apperrors.ErrInternal(err))
			c.Abort()
			return
		}

		store := orderrepo.NewMySQLOrder(db)
		productStore := productrepo.NewMySQLProduct(db)
		orderItemStore := orderitemrepo.NewMySQLOrder(db)
		userStore := userrepo.NewMySQLUser(db)
		authQueryClient := authtypes.NewQueryClient(appCtx.GRPCConn)
		bankClient := banktypes.NewQueryClient(appCtx.GRPCConn)
		cosmosStore := cosmosrepo.NewCosmos(appCtx.ClientCtx, appCtx.TxFactory, appCtx.Keyring, authQueryClient, bankClient, &appCtx.ProtoCodec)
		biz := orderbus.NewCreateOrder(store, productStore, orderItemStore, userStore, cosmosStore, db)
		recordID, err := biz.CreateOrder(c.Request.Context(), &data)

		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, appresponses.NewDataResponse(recordID, "create order successfully"))
	}

}
