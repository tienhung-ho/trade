package categorybusiness

import (
	"client/internal/common/apperrors"
	categorymodel "client/internal/model/mysql/category"
	"context"
	"errors"

	"gorm.io/gorm"
)

type GetCategoryStorage interface {
	GetCategory(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*categorymodel.Category, error)
}

type GetCategoryCache interface {
	GetCategory(ctx context.Context,
		cond map[string]interface{}, morekeys ...string) (*categorymodel.Category, error)
	SaveCategory(ctx context.Context, data interface{}, morekeys ...string) error
}

type getCategoryBusiness struct {
	store GetCategoryStorage
	cache GetCategoryCache
}

func NewGetCategoryBiz(store GetCategoryStorage, cache GetCategoryCache) *getCategoryBusiness {
	return &getCategoryBusiness{store, cache}
}

func (biz *getCategoryBusiness) GetCategory(ctx context.Context,
	cond map[string]interface{}, morekeys ...string) (*categorymodel.Category, error) {

	record, err := biz.store.GetCategory(ctx, cond, morekeys...)

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil,
				apperrors.ErrNotFoundEntity(categorymodel.EntityName, err)
		}

		return nil, apperrors.ErrCannotGetEntity(categorymodel.EntityName, err)
	}

	return record, nil
}
