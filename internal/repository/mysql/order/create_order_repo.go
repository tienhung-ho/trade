package orderrepo

import (
	"client/internal/common/apperrors"
	ordermodel "client/internal/model/mysql/order"
	responseutil "client/internal/util/response"
	"context"
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func (r *mysqlOrder) CreateOrder(ctx context.Context, db *gorm.DB,
	data *ordermodel.CreateOrder, morekeys ...string) (uint64, error) {

	if err := db.Omit("CreateOrderItems").Create(&data).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			fieldName := responseutil.ExtractFieldFromError(err, ordermodel.EntityName) // Extract field causing the duplicate error
			return 0, apperrors.ErrDuplicateEntry(ordermodel.EntityName, fieldName, err)
		}
		return 0, apperrors.ErrDB(err)
	}

	return data.OrderID, nil
}
