package userrepo

import (
	"client/internal/common/apperrors"
	usermodel "client/internal/model/mysql/user"
	responseutil "client/internal/util/response"
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func (s *mysqlUser) RegisterUser(db *gorm.DB,
	data *usermodel.UserRegister, morekeys ...string) (uint64, error) {

	if db.Error != nil {
		db.Rollback()
		return 0, apperrors.ErrDB(db.Error)
	}

	if err := db.Create(&data).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			fieldName := responseutil.ExtractFieldFromError(err, usermodel.EntityName) // Extract field causing the duplicate error
			db.Rollback()
			return 0, apperrors.ErrDuplicateEntry(usermodel.EntityName, fieldName, err)
		}
		db.Rollback()
		return 0, err
	}

	return data.ID, nil
}
