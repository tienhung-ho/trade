package productrepo

import (
	"client/internal/common/apperrors"
	productmodel "client/internal/model/mysql/product"
	responseutil "client/internal/util/response"
	"context"
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func (s *mysqlProduct) CreateProduct(ctx context.Context, db *gorm.DB,
	data *productmodel.CreateProduct, morekeys ...string) (uint64, error) {

	if db.Error != nil {
		return 0, apperrors.ErrDB(db.Error)
	}

	if err := db.WithContext(ctx).Create(&data).Error; err != nil {

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			fieldName := responseutil.ExtractFieldFromError(err, productmodel.EntityName) // Extract field causing the duplicate error
			return 0, apperrors.ErrDuplicateEntry(productmodel.EntityName, fieldName, err)
		}
		db.Rollback()
		return 0, err
	}

	return data.ProductID, nil
}

func (s *mysqlProduct) CreatePatchProducts(ctx context.Context,
	db *gorm.DB, data []productmodel.CreateProduct, morekeys ...string) error {

	if db.Error != nil {
		return apperrors.ErrDB(db.Error)
	}

	// Bắt đầu transaction nếu chưa có
	tx := db
	if tx.Statement == nil || tx.Statement.ConnPool == nil {
		tx = db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
	}

	// Tạo slice để lưu các sản phẩm cần tạo
	products := make([]productmodel.CreateProduct, len(data))
	copy(products, data)

	// Thực hiện tạo batch products
	if err := tx.WithContext(ctx).Create(&products).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			fieldName := responseutil.ExtractFieldFromError(err, productmodel.EntityName)
			tx.Rollback()
			return apperrors.ErrDuplicateEntry(productmodel.EntityName, fieldName, err)
		}
		tx.Rollback()
		return err
	}

	// Thu thập danh sách các ID đã được tạo
	productIDs := make([]uint64, len(products))
	for i, product := range products {
		productIDs[i] = product.ProductID
		// Cập nhật lại ID cho data ban đầu
		data[i].ProductID = product.ProductID
	}

	// Commit transaction nếu chúng ta đã bắt đầu một transaction mới
	if tx != db {
		if err := tx.Commit().Error; err != nil {
			return apperrors.ErrDB(err)
		}
	}

	return nil
}
