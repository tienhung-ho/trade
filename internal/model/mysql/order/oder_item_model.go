package ordermodel

import (
	productmodel "client/internal/model/mysql/product"
)

type OrderItem struct {
	OrderItemID uint64  `gorm:"column:order_item_id;primaryKey;autoIncrement" json:"order_item_id"`
	OrderID     uint64  `gorm:"column:order_id;not null" json:"order_id"`
	ProductID   uint64  `gorm:"column:product_id;not null" json:"product_id"`
	Quantity    uint    `gorm:"column:quantity;not null" json:"quantity"`
	UnitPrice   float64 `gorm:"column:unit_price;type:decimal(15,2);not null" json:"unit_price"`
	TotalPrice  float64 `gorm:"column:total_price;type:decimal(15,2);not null" json:"total_price"`
	Notes       string  `gorm:"column:notes;type:text" json:"notes,omitempty"`
	// Quan hệ
	//Order   Order                `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Product productmodel.Product `gorm:"foreignKey:ProductID;references:ProductID" json:"product,omitempty"`
}

func (OrderItem) TableName() string {
	return "order_item"
}
