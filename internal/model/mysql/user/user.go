package usermodel

import (
	"client/internal/common/datatypes"
	"client/internal/common/model"
	walletmodel "client/internal/model/mysql/wallet"
	"time"
)

// User đại diện cho bảng `user`
type User struct {
	ID        uint64         `gorm:"column:user_id;primaryKey;autoIncrement" json:"user_id"`
	Fullname  string         `gorm:"column:fullname;size:300" json:"fullname,omitempty"`
	Email     string         `gorm:"column:email;size:255;unique;not null" json:"email"`
	Password  string         `gorm:"column:password;size:255;not null" json:"-"` // không expose password trong JSON
	Phone     string         `gorm:"column:phone;size:20" json:"phone,omitempty"`
	Status    string         `gorm:"column:status;type:enum('Pending','Active','Inactive');default:'Pending'" json:"status"`
	Gender    string         `gorm:"column:gender;type:enum('Male','Female','Other')" json:"gender,omitempty"`
	Profile   datatypes.JSON `gorm:"column:profile;type:json" json:"profile,omitempty"`
	LastLogin *time.Time     `gorm:"column:last_login" json:"last_login,omitempty"`

	model.CommonFields
	// Quan hệ one-to-many với UserWallet
	Wallets walletmodel.UserWallet `gorm:"foreignKey:UserID;references:ID" json:"wallets,omitempty"`
}

func (User) TableName() string {
	return "user"
}
