package productv1

import (
	producthandler "client/api/handler/product"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

func ProductRouter(product *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	product.POST("/", producthandler.CreateProductHandler(db, rdb))
	product.GET("/:id", producthandler.GetProductHandler(db, rdb))
	product.PATCH("/:id", producthandler.UpdateProductHandler(db, rdb))
	product.DELETE("/:id", producthandler.DeleteProductHandler(db, rdb))
	product.GET("/list", producthandler.ListProductHandler(db, rdb))
}
