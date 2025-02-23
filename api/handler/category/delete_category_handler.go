package categoryhandler

import (
	"client/internal/common/apperrors"
	"client/internal/common/appresponses"
	categorystorage "client/internal/repository/mysql/category"
	categorybusiness "client/internal/service/category"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func DeleteCategoryHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, apperrors.ErrInternal(err))
			c.Abort()
			return
		}

		store := categorystorage.NewMySQLCategory(db)
		//cache := categorycache.NewRdbStorage(rdb)
		//image := imagestorage.NewMySQLImage(db)
		biz := categorybusiness.NewDeleteCategoryBiz(store, nil, nil)

		if err := biz.DeleteCategory(c, map[string]interface{}{"category_id": id}); err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, appresponses.NewDataResponse(true, "deleted category successfully"))

	}
}
