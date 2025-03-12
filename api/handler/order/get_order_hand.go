package orderhandler

import (
	"client/internal/common/apperrors"
	"client/internal/common/appresponses"
	orderrepo "client/internal/repository/mysql/order"
	orderbus "client/internal/service/order"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func GetOrderHandle(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, apperrors.ErrInternal(err))
			c.Abort()
			return
		}

		store := orderrepo.NewMySQLOrder(db)
		biz := orderbus.NewGetOrderBiz(store)

		record, err := biz.GetOrder(c.Request.Context(), map[string]interface{}{"order_id": id})

		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, appresponses.NewDataResponse(record, "get order successfully"))
	}
}
