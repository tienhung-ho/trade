package producthandler

import (
	"client/internal/common/apperrors"
	"client/internal/common/appresponses"
	productrepo "client/internal/repository/mysql/product"
	productbus "client/internal/service/product"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func GetProductHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, apperrors.ErrInternal(err))
			return
		}

		store := productrepo.NewMySQLProduct(db)

		biz := productbus.NewGetProductBiz(store)

		record, err := biz.GetProduct(c.Request.Context(), map[string]interface{}{"product_id": id})

		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, appresponses.NewDataResponse(record, "get product successfully!"))

	}

}
