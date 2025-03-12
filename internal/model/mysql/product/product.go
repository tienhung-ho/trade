package productmodel

import (
	"client/internal/common/model"
	categorymodel "client/internal/model/mysql/category"
	imagemodel "client/internal/model/mysql/image"
)

const (
	EntityName = "product"
)

type Product struct {
	ProductID   uint64  `gorm:"column:product_id;primaryKey;autoIncrement" json:"product_id,omitempty"`
	UserID      uint64  `gorm:"column:user_id;not null" json:"user_id,omitempty"` // FK đến User (seller)
	CategoryID  uint64  `gorm:"column:category_id" json:"category_id,omitempty"`  // FK đến Category, cho phép NULL
	Name        string  `gorm:"column:name;size:200;not null" json:"name,omitempty"`
	Description string  `gorm:"column:description;type:text" json:"description,omitempty"`
	Stock       int     `gorm:"column:stock;not null;default:0" json:"stock,omitempty"`
	Price       float64 `gorm:"column:price;type:decimal(15,2);not null;default:0" json:"price,omitempty"`
	model.CommonFields

	// Relations
	Category categorymodel.Category `gorm:"foreignKey:CategoryID;references:CategoryID" json:"category,omitempty"`
	Images   []imagemodel.Image     `gorm:"polymorphic:Resource;polymorphicValue:product" json:"images,omitempty"`
}

func (Product) TableName() string {
	return "product"
}
