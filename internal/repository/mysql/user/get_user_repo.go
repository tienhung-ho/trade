package userrepo

import (
	usermodel "client/internal/model/mysql/user"
	"context"

	"gorm.io/gorm"
)

func (s *mysqlUser) GetUser(ctx context.Context, cond map[string]interface{},
	morekeys ...string) (*usermodel.User, error) {

	db := s.db

	var record usermodel.User

	if err := db.WithContext(ctx).
		Select(SelectFields).
		Where(cond).
		Preload("Wallets", func(db *gorm.DB) *gorm.DB {
			return db.Select(WalletSelectField)
		}).
		First(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

func (s *mysqlUser) GetUserByWalletAddress(ctx context.Context, address string) (*usermodel.User, error) {
	db := s.db

	var user usermodel.User
	if err := db.Model(&usermodel.User{}).
		// Chọn các cột cần thiết
		Select(SelectFields).
		// Join sang bảng user_wallets
		Joins("JOIN user_wallets ON user.user_id = user_wallets.user_id").
		Where("user_wallets.wallet_address = ?", address).
		// preload ví (nếu muốn lấy chi tiết ví)
		Preload("Wallets", func(db *gorm.DB) *gorm.DB {
			return db.Select(WalletSelectField)
		}).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
