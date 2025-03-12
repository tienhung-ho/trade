package productrepo

import (
	"client/internal/common/apperrors"
	productmodel "client/internal/model/mysql/product"
	"context"

	"gorm.io/gorm"
)

func (r *mysqlProduct) DeleteProduct(ctx context.Context, db *gorm.DB,
	cond map[string]interface{}, morekyes ...string) error {

	if db.Error != nil {
		return apperrors.ErrDB(db.Error)
	}

	var product productmodel.Product

	if err := db.WithContext(ctx).Where(cond).Delete(&product).Error; err != nil {
		db.Rollback()
		return apperrors.ErrDB(err)
	}

	return nil
}
