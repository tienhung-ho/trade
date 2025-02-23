package categoryhandler

import (
	"client/internal/common/apperrors"
	"client/internal/common/appresponses"
	categorymodel "client/internal/model/mysql/category"
	categorystorage "client/internal/repository/mysql/category"
	imagestorage "client/internal/repository/mysql/image"
	categorybusiness "client/internal/service/category"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func UpdateCategoryHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, apperrors.ErrInternal(err))
			c.Abort()
			return
		}

		var data categorymodel.UpdateCategory

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, apperrors.ErrInternal(err))
			c.Abort()
			return
		}

		store := categorystorage.NewMySQLCategory(db)
		//cache := categorycache.NewRdbStorage(rdb)
		image := imagestorage.NewMySQLImage(db)
		biz := categorybusiness.NewUpdateCategoryBiz(store, nil, image, db)

		updatedCategory, err := biz.UpdateCategory(c, map[string]interface{}{"category_id": id}, &data)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, appresponses.NewDataResponse(updatedCategory, "update category successfully"))

	}
}
