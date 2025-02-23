package imagestorage

import (
	imagemodel "client/internal/model/mysql/image"
	"context"
)

func (s *mysqlImage) ListItem(ctx context.Context, cond map[string]interface{}) ([]imagemodel.Image, error) {

	db := s.db

	var records []imagemodel.Image
	if err := db.WithContext(ctx).Select(SelectFields).Where(cond).Find(&records).Error; err != nil {

		return nil, err
	}

	return records, nil
}
