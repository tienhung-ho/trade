package ordermodel

import (
	"client/internal/common/model"
	usermodel "client/internal/model/mysql/user"
)

const (
	EntityName     = "order"
	ItemEntityName = "order_item"
)

// Order struct đại diện cho đơn hàng
type Order struct {
	OrderID         uint64  `gorm:"column:order_id;primaryKey;autoIncrement" json:"order_id"`
	SellerID        uint64  `gorm:"column:seller_id;not null" json:"seller_id"`
	UserID          uint64  `gorm:"column:user_id;not null" json:"user_id"`
	TotalAmount     float64 `gorm:"column:total_amount;type:decimal(15,2);not null;default:0" json:"total_amount"`
	ShippingFee     float64 `gorm:"column:shipping_fee;type:decimal(10,2);not null;default:0" json:"shipping_fee"`
	DiscountAmount  float64 `gorm:"column:discount_amount;type:decimal(10,2);not null;default:0" json:"discount_amount"`
	FinalAmount     float64 `gorm:"column:final_amount;type:decimal(15,2);not null;default:0" json:"final_amount"`
	ShippingAddress string  `gorm:"column:shipping_address;type:text;not null" json:"shipping_address"`
	RecipientName   string  `gorm:"column:recipient_name;size:100;not null" json:"recipient_name"`
	RecipientPhone  string  `gorm:"column:recipient_phone;size:20;not null" json:"recipient_phone"`

	model.CommonFields

	// Quan hệ
	User       usermodel.User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OrderItems []OrderItem    `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
}

func (Order) TableName() string {
	return "order"
}

// CreateOrderResponse chứa danh sách orderIDs và mảng unsignedTxs dưới dạng base64
type CreateOrderResponse struct {
	OrderIDs    []uint64 `json:"order_ids"`
	UnsignedTxs []string `json:"unsigned_txs"`
}
