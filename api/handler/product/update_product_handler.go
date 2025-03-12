package producthandler

import (
	"client/internal/common/apperrors"
	"client/internal/common/appresponses"
	productmodel "client/internal/model/mysql/product"
	imagestorage "client/internal/repository/mysql/image"
	productrepo "client/internal/repository/mysql/product"
	productbus "client/internal/service/product"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func UpdateProductHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, apperrors.ErrInternal(err))
			c.Abort()
			return
		}

		var data productmodel.UpdateProduct

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusInternalServerError, apperrors.ErrInternal(err))
			c.Abort()
			return
		}

		store := productrepo.NewMySQLProduct(db)
		image := imagestorage.NewMySQLImage(db)
		biz := productbus.NewUpdateProductBiz(store, image, db)

		record, err := biz.UpdateProduct(c.Request.Context(), db, map[string]interface{}{"product_id": id}, &data)

		if err != nil {
			c.JSON(http.StatusInternalServerError, apperrors.ErrInternal(err))
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, appresponses.NewDataResponse(record, "updated data successfully"))
	}

}
