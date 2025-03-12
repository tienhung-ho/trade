package categorymodel

import (
	"client/internal/common/model"
	imagemodel "client/internal/model/mysql/image"
)

const (
	EntityName = "category"
)

type Category struct {
	CategoryID  uint64             `gorm:"column:category_id;primaryKey;autoIncrement" json:"category_id,omitempty"`
	Name        string             `gorm:"column:name;size:200;not null" json:"name,omitempty"`
	Description string             `gorm:"column:description;type:text" json:"description,omitempty"`
	Images      []imagemodel.Image `gorm:"polymorphic:Resource;polymorphicValue:category" json:"image,omitempty"`
	model.CommonFields
}

type CreateCategory struct {
	CategoryID  uint64             `gorm:"column:category_id;primaryKey;autoIncrement" json:"-"`
	Name        string             `gorm:"column:name;size:200;not null;unique" json:"name" validate:"required"`
	Description string             `gorm:"column:description;type:text" json:"description"`
	Images      []imagemodel.Image `gorm:"polymorphic:Resource;polymorphicValue:category" json:"images"`
	model.CommonFields
}

func (CreateCategory) TableName() string {
	return Category{}.TableName()
}

func (Category) TableName() string {
	return "category"
}

type UpdateCategory struct {
	CategoryID  uint64             `gorm:"column:category_id;primaryKey;autoIncrement" json:"-"`
	Name        string             `gorm:"column:name;size:200;not null;unique" json:"name" validate:"required"`
	Description string             `gorm:"column:description;type:text" json:"description"`
	Images      []imagemodel.Image `gorm:"polymorphic:Resource;polymorphicValue:category" json:"images"`

	model.CommonFields
}

func (UpdateCategory) TableName() string {
	return Category{}.TableName()
}
