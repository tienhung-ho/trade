package categoryhandler

import (
	"client/internal/common/apperrors"
	"client/internal/common/appresponses"
	categorymodel "client/internal/model/mysql/category"
	categorystorage "client/internal/repository/mysql/category"
	imagestorage "client/internal/repository/mysql/image"
	categorybusiness "client/internal/service/category"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
)

func CreateCategoryHandler(db *gorm.DB, rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		var data categorymodel.CreateCategory
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, apperrors.ErrInternal(err))
			c.Abort()
			return
		}
		validate := validator.New()

		err := validate.Struct(&data)
		if err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				//appErr := apperrors}.ErrValidation(validationErrors)
				c.JSON(http.StatusBadRequest, apperrors.ErrValidation(validationErrors))
				return
			}

			// Xử lý lỗi khác nếu có
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		store := categorystorage.NewMySQLCategory(db)
		cloud := imagestorage.NewMySQLImage(db)
		//cache := categorycache.NewRdbStorage(rdb)
		biz := categorybusiness.NewCreateCategoryBusiness(store, nil, cloud, db)

		recordId, err := biz.CreateCategory(c, &data)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, appresponses.NewDataResponse(recordId, "create category successfully"))
	}
}
