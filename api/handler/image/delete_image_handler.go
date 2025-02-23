package imagehandler

import (
	"client/internal/common/apperrors"
	"client/internal/common/appresponses"
	imagestorage "client/internal/repository/mysql/image"
	imagebusiness "client/internal/service/image"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func DeleteImageHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, apperrors.ErrInternal(err))
			c.Abort()
			return
		}

		store := imagestorage.NewMySQLImage(db)
		biz := imagebusiness.NewDeleteImageBiz(store)

		if err := biz.DeleteImage(c, map[string]interface{}{"image_id": id}); err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, appresponses.NewDataResponse(true, "delete image successfully"))
	}
}
