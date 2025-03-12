package productrepo

import (
	productmodel "client/internal/model/mysql/product"
	"context"
)

func (s *mysqlProduct) GetProduct(ctx context.Context,
	cond map[string]interface{}, morekeys ...string) (*productmodel.Product, error) {

	db := s.db

	var record productmodel.Product

	if err := db.WithContext(ctx).Where(cond).
		Preload("Images").
		Preload("Category").
		First(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}
