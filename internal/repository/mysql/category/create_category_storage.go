package categorystorage

import (
	"client/internal/common/apperrors"
	categorymodel "client/internal/model/mysql/category"
	responseutil "client/internal/util/response"
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func (s *mysqlCategory) CreateCategory(ctx context.Context,
	db *gorm.DB, data *categorymodel.CreateCategory, morekeys ...string) (uint64, error) {

	//db := s.db.Begin()

	if db.Error != nil {
		return 0, apperrors.ErrDB(db.Error)
	}

	if err := db.WithContext(ctx).Create(&data).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			fieldName := responseutil.ExtractFieldFromError(err, categorymodel.EntityName) // Extract field causing the duplicate error
			return 0, apperrors.ErrDuplicateEntry(categorymodel.EntityName, fieldName, err)
		}
		db.Rollback()
		return 0, err
	}

	//if err := db.Commit().Error; err != nil {
	//	db.Rollback()
	//	return 0, apperrors.ErrDB(err)
	//}

	return data.CategoryID, nil
}
