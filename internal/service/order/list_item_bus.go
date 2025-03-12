package orderbus

import (
	"client/internal/common/apperrors"
	"client/internal/common/filter"
	"client/internal/common/paging"
	ordermodel "client/internal/model/mysql/order"
	"context"
)

type ListOrderInterface interface {
	ListItem(ctx context.Context, cond map[string]interface{},
		paging *paging.Paging,
		filter *filter.Filter, morekeys ...string) ([]ordermodel.Order, error)
}

type ListOrderBusiness struct {
	store ListOrderInterface
}

func NewListOrderBiz(store ListOrderInterface) *ListOrderBusiness {
	return &ListOrderBusiness{
		store: store,
	}
}

func (biz *ListOrderBusiness) ListItem(ctx context.Context, cond map[string]interface{},
	paging *paging.Paging,
	filter *filter.Filter, morekeys ...string) ([]ordermodel.Order, error) {

	record, err := biz.store.ListItem(ctx, cond, paging, filter)

	if err != nil {
		return nil, apperrors.ErrCannotListEntity(ordermodel.EntityName, err)
	}

	return record, nil
}
