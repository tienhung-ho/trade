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
