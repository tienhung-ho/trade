package categorystorage

import (
	"client/internal/common/apperrors"
	"client/internal/common/filter"
	"client/internal/common/paging"
	categorymodel "client/internal/model/mysql/category"
	"context"

	"gorm.io/gorm"
)

func (s *mysqlCategory) ListItem(ctx context.Context,
	cond map[string]interface{}, paging *paging.Paging,
	filter *filter.Filter, morekeys ...string) ([]categorymodel.Category, error) {
	db := s.db

	// // Đếm tổng số lượng items
	if err := s.countRecord(db, cond, paging, filter); err != nil {
		return nil, err
	}

	// Xây dựng truy vấn động
	query := s.buildQuery(db, cond, filter)

	// Thêm phân trang
	query, err := s.addPaging(query, paging)
	if err != nil {
		return nil, err
	}

	// Thực hiện truy vấn
	var records []categorymodel.Category
	if err := query.Select(SelectFields).
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Select(ImageSelectFields)
		}).Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func (s *mysqlCategory) countRecord(db *gorm.DB, cond map[string]interface{}, paging *paging.Paging, filter *filter.Filter) error {
	// Apply conditions and filters
	db = s.buildQuery(db, cond, filter)

	if names, ok := cond["names"]; ok {
		db = db.Where("name IN ?", names)
	}

	if err := db.Table(categorymodel.Category{}.TableName()).Count(&paging.Total).Error; err != nil {
		return apperrors.NewErrorResponse(err, "Error count items from database", err.Error(), "CouldNotCount")
	}
	return nil
}

func (s *mysqlCategory) buildQuery(db *gorm.DB, cond map[string]interface{}, filter *filter.Filter) *gorm.DB {
	db = db.Where(cond)
	if filter != nil {
		if filter.Search != "" {
			searchPattern := "%" + filter.Search + "%"
			db = db.Where("name LIKE ?", searchPattern)
		}
		if filter.Status != "" {
			db = db.Where("status = ?", filter.Status)
		}

	}
	return db
}

func (s *mysqlCategory) addPaging(db *gorm.DB, paging *paging.Paging) (*gorm.DB, error) {
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
		db = db.Order("category_id desc")
	}

	// Apply pagination
	offset := (paging.Page - 1) * paging.Limit
	db = db.Offset(offset).Limit(paging.Limit)

	return db, nil
}
