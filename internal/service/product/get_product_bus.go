package productbus

import (
	"client/internal/common/apperrors"
	productmodel "client/internal/model/mysql/product"
	"context"
)

type GetProductInterface interface {
	GetProduct(ctx context.Context,
		cond map[string]interface{}, morekeys ...string) (*productmodel.Product, error)
}

type GetProductBusiness struct {
	store GetProductInterface
}

func NewGetProductBiz(store GetProductInterface) *GetProductBusiness {
	return &GetProductBusiness{
		store: store,
	}
}

func (biz *GetProductBusiness) GetProduct(ctx context.Context,
	cond map[string]interface{}, morekeys ...string) (*productmodel.Product, error) {

	record, err := biz.store.GetProduct(ctx, cond)

	if err != nil {
		return nil, apperrors.ErrCannotGetEntity(productmodel.EntityName, err)
	}

	return record, nil

}
