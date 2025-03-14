package walletmodel

import (
	"client/internal/common/datatypes"
	"client/internal/common/model"
)

const (
	EntityName = "wallet"
)

// UserWallet đại diện cho bảng `user_wallets`
type UserWallet struct {
	WalletID          uint64         `gorm:"column:wallet_id;primaryKey;autoIncrement" json:"wallet_id,omitempty"`
	UserID            uint64         `gorm:"column:user_id;not null" json:"user_id,omitempty"`
	WalletAddress     string         `gorm:"column:wallet_address;not null" json:"wallet_address,omitempty"`
	EncryptedMnemonic string         `gorm:"column:encrypted_mnemonic;type:text" json:"-"`
	WalletType        string         `gorm:"column:wallet_type;type:enum('ethereum','bitcoin','cosmos','solana','citcoin','other');default:'other'" json:"wallet_type,omitempty"`
	Balance           string         `gorm:"column:balance;type:decimal(65,18);default:0" json:"balance,omitempty"`
	Metadata          datatypes.JSON `gorm:"column:metadata;type:json" json:"metadata,omitempty"`

	model.CommonFields
	// Quan hệ many-to-one với User
	//User usermodel.User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

func (UserWallet) TableName() string {
	return "user_wallets"
}

func NewUserWallet(userID uint64, walletAddr, encryptedMnemonic,
	walletType, balance string, metadata datatypes.JSON) *UserWallet {

	return &UserWallet{
		UserID:            userID,
		WalletAddress:     walletAddr,
		EncryptedMnemonic: encryptedMnemonic,
		WalletType:        walletType,
		Balance:           balance,
		Metadata:          metadata,
	}
}
