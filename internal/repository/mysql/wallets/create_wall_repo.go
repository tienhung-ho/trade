package walletrepo

import (
	"client/internal/common/apperrors"
	walletmodel "client/internal/model/mysql/wallet"
	responseutil "client/internal/util/response"
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func (s *mysqlWallet) CreateWallet(db *gorm.DB,
	data *walletmodel.UserWallet, morekeys ...string) (uint64, uint64, error) {

	if err := db.Create(data).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			fieldName := responseutil.ExtractFieldFromError(err, walletmodel.EntityName) // Extract field causing the duplicate error
			db.Rollback()

			return 0, 0, apperrors.ErrDuplicateEntry(walletmodel.EntityName, fieldName, err)
		}
		db.Rollback()

		return 0, 0, err
	}

	return data.WalletID, data.UserID, nil
}
