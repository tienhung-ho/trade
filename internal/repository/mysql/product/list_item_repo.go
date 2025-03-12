package productrepo

import (
	"client/internal/common/apperrors"
	"client/internal/common/filter"
	"client/internal/common/paging"
	productmodel "client/internal/model/mysql/product"
	"context"

	"gorm.io/gorm"
)

func (r *mysqlProduct) ListItem(ctx context.Context, cond map[string]interface{},
	paging *paging.Paging, filter *filter.Filter, morekeys ...string) ([]productmodel.Product, error) {
	db := r.db.WithContext(ctx)

	// Count total records
	if err := r.countRecord(db, cond, paging, filter); err != nil {
		return nil, err
	}

	// Build query with filters
	query := r.buildQuery(db, cond, filter)

	// Add pagination and sorting
	query, err := r.addPaging(query, paging)
	if err != nil {
		return nil, err
	}
	// Execute query
	var products []productmodel.Product
	if err := query.
		Preload("Category").
		Preload("Images").
		Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *mysqlProduct) countRecord(db *gorm.DB, cond map[string]interface{},
	paging *paging.Paging, filter *filter.Filter) error {
	// Apply conditions and filters
	db = r.buildQuery(db, cond, filter)

	if err := db.Model(&productmodel.Product{}).Count(&paging.Total).Error; err != nil {
		return apperrors.NewErrorResponse(err, "Error counting items from database", err.Error(), "CouldNotCount")
	}
	return nil
}

func (r *mysqlProduct) buildQuery(db *gorm.DB, cond map[string]interface{}, filter *filter.Filter) *gorm.DB {
	db = db.Where(cond)
	if filter != nil {

		if len(filter.IDs) > 0 {
			db = db.Where("product_id IN (?)", filter.IDs)
		}

		if filter.Status != "" {
			db = db.Where("status = ?", filter.Status)
		}

		if filter.Search != "" {
			searchPattern := "%" + filter.Search + "%"
			db = db.Where("name LIKE ? OR description LIKE ?", searchPattern, searchPattern)
		}

		if filter.CategoryID != 0 {
			db = db.Where("category_id = ?", filter.CategoryID)
		}

	}
	return db
}

func (r *mysqlProduct) addPaging(db *gorm.DB, paging *paging.Paging) (*gorm.DB, error) {
	// Parse and validate the sort fields
	sortFields, err := paging.ParseSortFields(paging.Sort, AllowedSortFields)
	if err != nil {
		return nil, apperrors.NewErrorResponse(err, "Invalid sort parameters", err.Error(), "InvalidSort")
	}

	// Apply sorting to the query
	if len(sortFields) > 0 {
		for _, sortField := range sortFields {
			db = db.Order(sortField)
		}
	} else {
		// Default sorting if no sort parameters are provided
		db = db.Order("product_id asc")
	}

	// Apply pagination
	offset := (paging.Page - 1) * paging.Limit
	db = db.Offset(offset).Limit(paging.Limit)

	return db, nil
}

func (r *mysqlProduct) ListItemByIDs(ctx context.Context, cond []uint64) ([]productmodel.Product, error) {
	db := r.db.WithContext(ctx)

	var products []productmodel.Product
	if err := db.Where("product_id IN (?)", cond).
		Select("product_id, price, stock, user_id").
		Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}
