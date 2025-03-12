package producthandler

import (
	"client/internal/common/apperrors"
	"client/internal/common/appresponses"
	productmodel "client/internal/model/mysql/product"
	imagestorage "client/internal/repository/mysql/image"
	productrepo "client/internal/repository/mysql/product"
	productbus "client/internal/service/product"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func CreateProductHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {

	return func(c *gin.Context) {
		var data productmodel.CreateProduct

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusInternalServerError, apperrors.ErrInternal(err))
			c.Abort()
			return
		}

		store := productrepo.NewMySQLProduct(db)
		image := imagestorage.NewMySQLImage(db)
		biz := productbus.NewCreateProductBiz(store, image, db)

		record, err := biz.CreateProduct(c.Request.Context(), db, &data)
		if err != nil {

			c.JSON(http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, appresponses.NewDataResponse(record, "created new product successfully!"))

	}

}
