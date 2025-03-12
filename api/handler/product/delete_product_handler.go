package producthandler

import (
	"client/internal/common/apperrors"
	"client/internal/common/appresponses"
	imagestorage "client/internal/repository/mysql/image"
	productrepo "client/internal/repository/mysql/product"
	productbus "client/internal/service/product"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func DeleteProductHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, apperrors.ErrInternal(err))
			return
		}

		store := productrepo.NewMySQLProduct(db)
		image := imagestorage.NewMySQLImage(db)

		biz := productbus.NewDeleteProducBiz(store, image, db)

		if err := biz.DeleteProduct(c.Request.Context(), db, map[string]interface{}{"product_id": id}); err != nil {
			c.JSON(http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, appresponses.NewDataResponse(true, "deleted product successfully!"))
	}

}
