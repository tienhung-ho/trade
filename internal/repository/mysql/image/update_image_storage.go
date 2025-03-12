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

func (s *mysqlImage) BulkUpdateResourceID(ctx context.Context, db *gorm.DB,
	imageIDs []uint64, resourceID *uint64, resourceType *string) error {
	if len(imageIDs) == 0 {
		return nil
	}

	tx := db.WithContext(ctx).Model(&imagemodel.Image{}).Where("image_id IN ?", imageIDs)

	updateData := map[string]interface{}{}

	if resourceID == nil {
		updateData["resource_id"] = 0
		updateData["resource_type"] = ""
	} else {
		updateData["resource_id"] = *resourceID
		if resourceType != nil {
			updateData["resource_type"] = *resourceType
		}
	}

	return tx.Updates(updateData).Error
}

// BulkDeleteImages: xoá hoàn toàn record image cho 1 list ID
func (s *mysqlImage) BulkDeleteImages(ctx context.Context, db *gorm.DB, imageIDs []uint64) error {
	if len(imageIDs) == 0 {
		return nil
	}
	return db.WithContext(ctx).Where("image_id IN ?", imageIDs).
		Delete(&imagemodel.UpdateImage{}).Error
}
