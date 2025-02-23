package imagestorage

import (
	"client/internal/common/apperrors"
	imagemodel "client/internal/model/mysql/image"
	"context"
)

func (s *mysqlImage) DeleteImage(ctx context.Context, cond map[string]interface{}, morekeys ...string) error {

	db := s.db.Begin()

	if db.Error != nil {
		return apperrors.ErrDB(db.Error)
	}

	var image imagemodel.Image
	if err := db.WithContext(ctx).Where(cond).Delete(&image).Error; err != nil {
		db.Rollback()
		return err
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return apperrors.ErrDB(err)
	}

	return nil
}
