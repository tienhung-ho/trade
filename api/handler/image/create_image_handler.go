package imagehandler

import (
	"client/internal/common/apperrors"
	"client/internal/common/appresponses"
	imagestorage "client/internal/repository/mysql/image"
	imagebusiness "client/internal/service/image"
	cloudutil "client/internal/util/cloud"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
)

func CreateImageHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Lấy dữ liệu từ context mà middleware đã đặt
		fileBufferInterface, exists := c.Get("fileBuffer")
		if !exists {
			c.JSON(http.StatusInternalServerError, apperrors.ErrInternal(nil))
			return
		}
		fileNameInterface, exists := c.Get("fileName")
		if !exists {
			c.JSON(http.StatusInternalServerError, apperrors.ErrInternal(nil))
			return
		}

		fileBuffer, ok := fileBufferInterface.([]byte)
		if !ok {
			c.JSON(http.StatusInternalServerError, apperrors.ErrInternal(nil))
			return
		}
		fileName, ok := fileNameInterface.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, apperrors.ErrInternal(nil))
			return
		}

		uniqueFileName := uuid.New().String() + "-" + fileName
		// Chuẩn bị dữ liệu để gửi vào business layer
		data := &cloudutil.Image{
			FileName:   uniqueFileName,
			FileBuffer: fileBuffer,
			// Thêm các trường khác nếu cần
		}

		// Gọi business layer để xử lý
		store := imagestorage.NewMySQLImage(db)
		biz := imagebusiness.NewCreateImageBiz(store)
		recordID, fileURL, err := biz.CreateImage(c.Request.Context(), data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		// Trả về phản hồi cho client
		resp := struct {
			ImageID  uint64 `json:"image_id"`
			ImageURL string `json:"image_url"`
			Message  string `json:"message"`
		}{
			ImageID:  recordID,
			ImageURL: fileURL,
			Message:  "Image created successfully",
		}

		c.JSON(http.StatusOK, appresponses.NewDataResponse(resp, "create image successfully"))
	}
}
