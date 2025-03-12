package orderitemrepo

import (
	"client/internal/common/apperrors"
	ordermodel "client/internal/model/mysql/order"
	responseutil "client/internal/util/response"
	"context"
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func (s *mysqlOrderItem) CreatePatchOrderItem(ctx context.Context,
	db *gorm.DB, data []ordermodel.CreateOrderItem, morekeys ...string) error {

	// Thực hiện bulk insert
	if err := db.WithContext(ctx).Create(&data).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				// Lỗi trùng lặp
				fieldName := responseutil.ExtractFieldFromError(err, ordermodel.ItemEntityName)
				return apperrors.ErrDuplicateEntry(ordermodel.ItemEntityName, fieldName, err)
			case 1452:
				// Lỗi khóa ngoại
				return apperrors.ErrForeignKeyConstraint(ordermodel.ItemEntityName, "", err) // Sửa "ingredient_id" thành "recipe_id"
			default:
				// Các lỗi MySQL khác
				return apperrors.ErrDB(err)
			}
		}
		return apperrors.ErrDB(err)
	}

	// Không gọi tx.Commit() hay tx.Rollback() ở đây

	return nil
}
