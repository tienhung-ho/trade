package orderv1

import (
	orderhandler "client/api/handler/order"
	cosmosmodel "client/internal/model/cosmos"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func OrderRouter(order *gin.RouterGroup, db *gorm.DB, rdb *redis.Client, cosmos *cosmosmodel.AppContext) {
	order.POST("/", orderhandler.CreateOrderhandler(db, rdb, cosmos))
	order.GET("/:id", orderhandler.GetOrderHandle(db, rdb))
	order.GET("/list", orderhandler.ListOrderHand(db, rdb))
}
