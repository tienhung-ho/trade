package categoryv1

import (
	categoryhandler "client/api/handler/category"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

func CategoryRouter(category *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	category.GET("/:id", categoryhandler.GetCategoryHandler(db, rdb))
	category.POST("/", categoryhandler.CreateCategoryHandler(db, rdb))
	category.PATCH("/:id", categoryhandler.UpdateCategoryHandler(db, rdb))
	category.DELETE("/:id", categoryhandler.DeleteCategoryHandler(db, rdb))
	category.GET("/list", categoryhandler.ListCategoryHandler(db, rdb))
}
