package imagestorage

import (
	"client/internal/common/apperrors"
	imagemodel "client/internal/model/mysql/image"
	"context"
	"gorm.io/gorm"
)

func (s *mysqlImage) UpdateImage(ctx context.Context, cond map[string]interface{}, data *imagemodel.UpdateImage) error {

	db := s.db.Begin()

	if db.Error != nil {
		return apperrors.ErrDB(db.Error)
	}

	if err := db.WithContext(ctx).Where(cond).Updates(data).Error; err != nil {
		db.Rollback()
		return err
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return apperrors.ErrDB(err)
	}

	return nil
}

// BulkUpdateResourceID: cập nhật resource_id cho 1 list image_id
func (s *mysqlImage) BulkUpdateResourceID(ctx context.Context, db *gorm.DB, imageIDs []uint64, resourceID *uint64) error {
	if len(imageIDs) == 0 {
		return nil
	}
	tx := db.WithContext(ctx)
	if resourceID == nil {
		return tx.Model(&imagemodel.UpdateImage{}).Where("image_id IN ?", imageIDs).
			Update("resource_id", nil).Error
	}
	return tx.Model(&imagemodel.UpdateImage{}).Where("image_id IN ?", imageIDs).
		Update("resource_id", *resourceID).Error
}

// BulkDeleteImages: xoá hoàn toàn record image cho 1 list ID
func (s *mysqlImage) BulkDeleteImages(ctx context.Context, db *gorm.DB, imageIDs []uint64) error {
	if len(imageIDs) == 0 {
		return nil
	}
	return db.WithContext(ctx).Where("image_id IN ?", imageIDs).
		Delete(&imagemodel.UpdateImage{}).Error
}
