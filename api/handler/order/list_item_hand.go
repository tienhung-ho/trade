package orderhandler

import (
	"client/internal/common/apperrors"
	"client/internal/common/appresponses"
	"client/internal/common/filter"
	"client/internal/common/paging"
	orderrepo "client/internal/repository/mysql/order"
	orderbus "client/internal/service/order"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func ListOrderHand(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {

	return func(c *gin.Context) {

		condition := map[string]interface{}{}

		var paging paging.Paging

		if err := c.ShouldBindJSON(&paging); err != nil {
			c.JSON(http.StatusBadRequest, apperrors.ErrInternal(err))
			c.Abort()
			return
		}

		var filter filter.Filter

		if err := c.ShouldBindJSON(&filter); err != nil {
			c.JSON(http.StatusBadRequest, apperrors.ErrInternal(err))
			c.Abort()
			return
		}

		store := orderrepo.NewMySQLOrder(db)
		biz := orderbus.NewListOrderBiz(store)

		record, err := biz.ListItem(c.Request.Context(), condition, &paging, &filter)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, appresponses.NewSuccesResponse(record, paging, filter))
	}

}
