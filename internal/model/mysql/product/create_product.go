package productmodel

import (
	categorymodel "client/internal/model/mysql/category"
	imagemodel "client/internal/model/mysql/image"
	usermodel "client/internal/model/mysql/user"
)

type CreateProduct struct {
	ProductID   uint64  `gorm:"column:product_id;primaryKey;autoIncrement" json:"product_id"`
	UserID      uint64  `gorm:"column:user_id;not null" json:"user_id"`          // FK đến User (seller)
	CategoryID  uint64  `gorm:"column:category_id" json:"category_id,omitempty"` // FK đến Category, cho phép NULL
	Name        string  `gorm:"column:name;size:200;not null" json:"name"`
	Description string  `gorm:"column:description;type:text" json:"description,omitempty"`
	Stock       int     `gorm:"column:stock;not null;default:0" json:"stock"`
	Price       float64 `gorm:"column:price;type:decimal(15,2);not null;default:0" json:"price"`

	// Relations
	Images   []imagemodel.Image     `gorm:"polymorphic:Resource;polymorphicValue:product" json:"images"`
	User     usermodel.User         `gorm:"-" json:"user,omitempty"`
	Category categorymodel.Category `gorm:"-" json:"category,omitempty"`
}

func (CreateProduct) TableName() string {
	return Product{}.TableName()
}
