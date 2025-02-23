package categorystorage

import (
	categorymodel "client/internal/model/mysql/category"
	"context"
)

func (s *mysqlCategory) DeleteCategory(ctx context.Context, cond map[string]interface{}, morekeys ...string) error {

	db := s.db

	var record categorymodel.Category
	if err := db.WithContext(ctx).Where(cond).Delete(&record).Error; err != nil {
		return err
	}

	return nil
}
