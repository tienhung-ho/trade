package categorystorage

import (
	categorymodel "client/internal/model/mysql/category"
	"context"

	"gorm.io/gorm"
)

func (s *mysqlCategory) GetCategory(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*categorymodel.Category, error) {

	db := s.db

	var record categorymodel.Category
	if err := db.WithContext(ctx).
		Where(cond).Preload("Images", func(db *gorm.DB) *gorm.DB {
		return db.Select(ImageSelectFields)
	}).
		First(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}
