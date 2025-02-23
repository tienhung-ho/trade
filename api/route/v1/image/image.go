package imagev1

import (
	imagehandler "client/api/handler/image"
	cloudmiddleware "client/api/middleware/cloud"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ImageRouter(image *gin.RouterGroup, db *gorm.DB) {
	image.POST("/", cloudmiddleware.ImageValidationMiddleware(), imagehandler.CreateImageHandler(db))
	image.DELETE("/:id", imagehandler.DeleteImageHandler(db))
}
