package ordermodel

type CreateOrderItem struct {
	OrderItemID uint64  `gorm:"column:order_item_id;primaryKey;autoIncrement" json:"order_item_id"`
	OrderID     uint64  `gorm:"column:order_id;not null"`
	ProductID   uint64  `gorm:"column:product_id;not null" json:"product_id"`
	Quantity    uint    `gorm:"column:quantity;not null" json:"quantity"`
	UnitPrice   float64 `gorm:"column:unit_price;type:decimal(15,2);not null"`
	TotalPrice  float64 `gorm:"column:total_price;type:decimal(15,2);not null"`
	Notes       string  `gorm:"column:notes;type:text" json:"notes,omitempty"`
	// Quan hệ
	//	Order Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

func (CreateOrderItem) TableName() string {
	return OrderItem{}.TableName()
}
