package productmodel

import (
	"client/internal/common/model"
	categorymodel "client/internal/model/mysql/category"
	imagemodel "client/internal/model/mysql/image"
)

type UpdateProduct struct {
	ProductID   uint64  `gorm:"column:product_id;primaryKey;autoIncrement" json:"product_id"`
	UserID      uint64  `gorm:"column:user_id;not null" json:"user_id"`          // FK đến User (seller)
	CategoryID  uint64  `gorm:"column:category_id" json:"category_id,omitempty"` // FK đến Category, cho phép NULL
	Name        string  `gorm:"column:name;size:200;not null" json:"name"`
	Description string  `gorm:"column:description;type:text" json:"description,omitempty"`
	Stock       int     `gorm:"column:stock;not null;default:0" json:"stock"`
	Price       float64 `gorm:"column:price;type:decimal(15,2);not null;default:0" json:"price"`
	model.CommonFields

	// Relations
	Category categorymodel.Category `gorm:"foreignKey:CategoryID;references:CategoryID" json:"category,omitempty"`
	Images   []imagemodel.Image     `gorm:"polymorphic:Resource;polymorphicValue:product" json:"images"`
}

func (UpdateProduct) TableName() string {
	return Product{}.TableName()
}

type ProductQuantityUpdate struct {
	ProductID  uint64
	Adjustment int
}

func (ProductQuantityUpdate) TableName() string {
	return Product{}.TableName()
}
