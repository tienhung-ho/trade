package orderrepo

import (
	"client/internal/common/apperrors"
	ordermodel "client/internal/model/mysql/order"
	"context"
	"errors"

	"gorm.io/gorm"
)

func (r *mysqlOrder) GetOrder(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*ordermodel.Order, error) {

	var order ordermodel.Order

	if err := r.db.WithContext(ctx).
		Where(cond).
		Preload("OrderItems").
		Preload("OrderItems.Product").
		First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.RecordNotFound
		}
		return nil, err
	}

	return &order, nil
}
