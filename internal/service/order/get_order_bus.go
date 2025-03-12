package orderbus

import (
	"client/internal/common/apperrors"
	ordermodel "client/internal/model/mysql/order"
	"context"
)

type GetOrderInterface interface {
	GetOrder(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*ordermodel.Order, error)
}

type GetOrderBusiness struct {
	store GetOrderInterface
}

func NewGetOrderBiz(store GetOrderInterface) *GetOrderBusiness {
	return &GetOrderBusiness{
		store: store,
	}
}

func (biz *GetOrderBusiness) GetOrder(ctx context.Context,
	cond map[string]interface{}, morekeys ...string) (*ordermodel.Order, error) {

	record, err := biz.store.GetOrder(ctx, cond)

	if err != nil {
		if err == apperrors.RecordNotFound {
			return nil, apperrors.ErrNotFoundEntity(ordermodel.EntityName, err)
		}

		return nil, apperrors.ErrCannotGetEntity(ordermodel.EntityName, err)
	}

	return record, nil
}
