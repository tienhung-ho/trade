package categorystorage

import (
	"client/internal/common/apperrors"
	categorymodel "client/internal/model/mysql/category"
	responseutil "client/internal/util/response"
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *mysqlCategory) UpdateCategory(ctx context.Context,
	db *gorm.DB,
	cond map[string]interface{},
	data *categorymodel.UpdateCategory,
	morekeys ...string) (*categorymodel.Category, error) {

	if db.Error != nil {
		return nil, apperrors.ErrDB(db.Error)
	}

	// Sử dụng con trỏ đến struct thực tế
	if err := db.WithContext(ctx).Model(&categorymodel.UpdateCategory{}).Where(cond).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Updates(data).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			fieldName := responseutil.ExtractFieldFromError(err, categorymodel.EntityName) // Extract field causing the duplicate error
			return nil, apperrors.ErrDuplicateEntry(categorymodel.EntityName, fieldName, err)
		}
		db.Rollback()
		return nil, err
	}

	var updatedCategory categorymodel.Category
	if err := db.WithContext(ctx).
		Select(SelectFields).
		Where(cond).
		Preload("Images").
		First(&updatedCategory).Error; err != nil {
		db.Rollback()
		return nil, apperrors.ErrNotFoundEntity(categorymodel.EntityName, err)
	}

	return &updatedCategory, nil
}
