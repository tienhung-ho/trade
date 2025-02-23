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

func GetCategoryHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, apperrors.ErrInternal(err))
			c.Abort()
			return
		}

		store := categorystorage.NewMySQLCategory(db)
		//cache := categorycache.NewRdbStorage(rdb)
		biz := categorybusiness.NewGetCategoryBiz(store, nil)

		record, err := biz.GetCategory(c.Request.Context(), map[string]interface{}{"category_id": id})

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, appresponses.NewDataResponse(record, "get account successfully"))

	}
}
