package categorybusiness

import (
	"client/internal/common/apperrors"
	"client/internal/common/filter"
	"client/internal/common/paging"
	categorymodel "client/internal/model/mysql/category"
	"context"
	"errors"
	
	"gorm.io/gorm"
)

type ListItemSotrage interface {
	ListItem(ctx context.Context, cond map[string]interface{},
		paging *paging.Paging, filter *filter.Filter, morekeys ...string) ([]categorymodel.Category, error)
}

type ListItemCache interface {
	ListItem(ctx context.Context, key string) ([]categorymodel.Category, error)
	SaveCategory(ctx context.Context, data interface{}, morekeys ...string) error
	SavePaging(ctx context.Context, paging *paging.Paging, morekeys ...string) error
	SaveFilter(ctx context.Context, filter *filter.Filter, morekeys ...string) error
	GetPaging(ctx context.Context, key string) (*paging.Paging, error)
}

type listItemCategoryBusiness struct {
	store ListItemSotrage
	cache ListItemCache
}

func NewListItemCategoryBiz(store ListItemSotrage, cache ListItemCache) *listItemCategoryBusiness {
	return &listItemCategoryBusiness{store, cache}
}

func (biz *listItemCategoryBusiness) ListItem(ctx context.Context,
	cond map[string]interface{}, paging *paging.Paging,
	filter *filter.Filter, morekeys ...string) ([]categorymodel.Category, error) {

	records, err := biz.store.ListItem(ctx, cond, paging, filter, morekeys...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, apperrors.ErrNotFoundEntity(categorymodel.EntityName, err)
		}

		return nil, apperrors.ErrCannotListEntity(categorymodel.EntityName, err)
	}
	return records, nil
}
