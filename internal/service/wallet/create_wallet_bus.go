package walletbus

import (
	"client/internal/common/apperrors"
	walletmodel "client/internal/model/mysql/wallet"

	"gorm.io/gorm"
)

type WalletInterface interface {
	CreateWallet(db *gorm.DB,
		data *walletmodel.UserWallet, morekeys ...string) (uint64, uint64, error)
}

type WalletBusiness struct {
	store WalletInterface
	db    *gorm.DB
}

func NewWalletBiz(store WalletInterface, db *gorm.DB) *WalletBusiness {
	return &WalletBusiness{
		store: store,
		db:    db,
	}
}

func (biz *WalletBusiness) CreateWallet(data *walletmodel.UserWallet, morekeys ...string) (uint64, error) {

	tx := biz.db.Begin()
	if err := tx.Error; err != nil {
		return 0, apperrors.ErrDB(err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	walletID, _, err := biz.store.CreateWallet(tx, data)

	if err != nil {
		tx.Rollback()
		return 0, apperrors.ErrCannotCreateEntity(walletmodel.EntityName, err)
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		//		_ = biz.cosmosStore.DeleteKey(alias)
		return 0, apperrors.ErrDB(err)
	}

	return walletID, nil
}
