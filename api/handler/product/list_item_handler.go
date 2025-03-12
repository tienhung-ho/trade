package producthandler

import (
	"client/internal/common/apperrors"
	"client/internal/common/appresponses"
	"client/internal/common/filter"
	"client/internal/common/paging"
	productmodel "client/internal/model/mysql/product"
	productrepo "client/internal/repository/mysql/product"
	productbus "client/internal/service/product"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func ListProductHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
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

		store := productrepo.NewMySQLProduct(db)
		biz := productbus.NewListItemBiz(store)
		records, err := biz.ListItem(c.Request.Context(), condition, &paging, &filter)

		if err != nil {
			fmt.Printf("Error sorting products: %v\n", err)
			// Return the error to the client
			var appErr *apperrors.AppError
			if errors.As(err, &appErr) {
				c.JSON(appErr.StatusCode, apperrors.ErrCannotSort(productmodel.EntityName, err))
			}
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, appresponses.NewSuccesResponse(records, paging, filter))

	}
}
