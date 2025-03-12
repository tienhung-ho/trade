package productrepo

import (
	"client/internal/common/apperrors"
	productmodel "client/internal/model/mysql/product"
	responseutil "client/internal/util/response"
	"context"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *mysqlProduct) UpdateProduct(ctx context.Context, db *gorm.DB, cond map[string]interface{},
	data *productmodel.UpdateProduct, morekeys ...string) (*productmodel.Product, error) {

	if db.Error != nil {
		return nil, apperrors.ErrDB(db.Error)
	}

	if err := db.WithContext(ctx).Model(&productmodel.UpdateProduct{}).Where(cond).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		//Clauses(clause.Returning{}).
		Updates(data).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			fieldName := responseutil.ExtractFieldFromError(err, productmodel.EntityName) // Extract field causing the duplicate error
			return nil, apperrors.ErrDuplicateEntry(productmodel.EntityName, fieldName, err)
		}
		db.Rollback()
		return nil, err
	}

	var record productmodel.Product

	if err := db.WithContext(ctx).Model(data).
		Where(cond).Preload("Images").
		Preload("Category").
		First(&record).Error; err != nil {

		db.Rollback()
		return nil, apperrors.ErrDB(err)
	}

	return &record, nil
}

func (r *mysqlProduct) BulkUpdateProductQuantity(ctx context.Context, db *gorm.DB,
	updates []productmodel.ProductQuantityUpdate) error {

	if db == nil {
		db = r.db
	}

	if db.Error != nil {
		return apperrors.ErrDB(db.Error)
	}

	// Start transaction if not already in one
	tx := db

	// Validate first - collect all product IDs
	var productIDs []uint64
	adjustmentMap := make(map[uint64]int)

	for _, update := range updates {
		if update.ProductID == 0 {
			tx.Rollback()
			return apperrors.ErrInvalidRequest(errors.New("product ID cannot be 0"))
		}
		productIDs = append(productIDs, update.ProductID)
		adjustmentMap[update.ProductID] = update.Adjustment
	}

	// Fetch all products that need to be updated in a single query
	var products []productmodel.Product
	if err := tx.WithContext(ctx).
		Table("product").
		Where("product_id IN ?", productIDs).
		Find(&products).Error; err != nil {
		tx.Rollback()
		return apperrors.ErrDB(err)
	}

	// Validate stock won't go negative
	for _, product := range products {
		adjustment := adjustmentMap[product.ProductID]
		newStock := product.Stock + adjustment
		if newStock < 0 {
			tx.Rollback()
			return apperrors.NewErrorResponse(
				nil,
				fmt.Sprintf("product %s (ID: %d) would have negative stock (%d current, %d adjustment)",
					product.Name, product.ProductID, product.Stock, adjustment),
				"Product",
				"InsufficientStock",
			)
		}
	}

	// Use raw SQL for bulk update with CASE statement
	// This is much more efficient than individual updates
	if len(updates) > 0 {
		query := "UPDATE product SET stock = CASE product_id "
		var params []interface{}

		for _, product := range products {
			adjustment := adjustmentMap[product.ProductID]
			query += "WHEN ? THEN stock + ? "
			params = append(params, product.ProductID, adjustment)
		}

		query += "ELSE stock END WHERE product_id IN ?"
		params = append(params, productIDs)

		if err := tx.WithContext(ctx).Table("product").Exec(query, params...).Error; err != nil {
			tx.Rollback()
			return apperrors.ErrDB(err)
		}
	}

	return nil
}
