package categoryhandler

import (
	"client/internal/common/apperrors"
	"client/internal/common/appresponses"
	"client/internal/common/filter"
	paging "client/internal/common/paging"
	categorystorage "client/internal/repository/mysql/category"
	categorybusiness "client/internal/service/category"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
)

func ListCategoryHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		condition := map[string]interface{}{
			//"status": []string{"pending", "active", "inactive"},
		}

		var paging paging.Paging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, apperrors.ErrInternal(err))
			return
		}

		paging.Process()

		var filter filter.Filter

		if err := c.ShouldBind(&filter); err != nil {
			c.JSON(http.StatusBadRequest, apperrors.ErrInternal(err))
			return
		}

		store := categorystorage.NewMySQLCategory(db)
		//cache := categorycache.NewRdbStorage(rdb)
		biz := categorybusiness.NewListItemCategoryBiz(store, nil)

		records, err := biz.ListItem(c.Request.Context(), condition, &paging, &filter)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, appresponses.NewSuccesResponse(records, paging, filter))
	}
}
