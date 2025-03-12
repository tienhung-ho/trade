package productbus

import (
	"client/internal/common/apperrors"
	"client/internal/common/filter"
	"client/internal/common/paging"
	productmodel "client/internal/model/mysql/product"
	"context"
)

type ListItemInterface interface {
	ListItem(ctx context.Context, cond map[string]interface{},
		paging *paging.Paging, filter *filter.Filter, morekeys ...string) ([]productmodel.Product, error)
}

type ListItemBusiness struct {
	store ListItemInterface
}

func NewListItemBiz(store ListItemInterface) *ListItemBusiness {
	return &ListItemBusiness{
		store: store,
	}
}

func (biz *ListItemBusiness) ListItem(ctx context.Context, cond map[string]interface{},
	paging *paging.Paging, filter *filter.Filter, morekeys ...string) ([]productmodel.Product, error) {

	record, err := biz.store.ListItem(ctx, cond, paging, filter)

	if err != nil {
		return nil, apperrors.ErrCannotListEntity(productmodel.EntityName, err)
	}

	return record, nil
}
