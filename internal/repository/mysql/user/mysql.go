package userrepo

import (
	"client/internal/common/apperrors"
	"client/internal/common/transaction"
	"context"

	"gorm.io/gorm"
)

var (
	AllowedSortFields = map[string]bool{}
	SelectFields      = []string{
		"user_id",  // AccountID
		"phone",    // Phone
		"fullname", // Fullname
		"status",   // Status
		"email",    // Email
		"gender",   // Gender
		"profile",
		"password",
	}
	WalletSelectField = []string{
		"wallet_id",
		"user_id",
		"wallet_address",
		"encrypted_mnemonic",
		"wallet_type",
		"balance",
	}
)

type mysqlUser struct {
	db *gorm.DB
}

func NewMySQLOrder(db *gorm.DB) *mysqlUser {
	return &mysqlUser{db}
}

func (r *mysqlUser) Transaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return apperrors.ErrDB(tx.Error)
	}

	txCtx := context.WithValue(ctx, transaction.TransactionKey, tx)

	if err := fn(txCtx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return apperrors.ErrDB(err)
	}

	return nil
}

// getDB lấy *gorm.DB từ context nếu có transaction
func (r *mysqlUser) getDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(transaction.TransactionKey).(*gorm.DB)
	if ok {
		return tx
	}
	return r.db
}
