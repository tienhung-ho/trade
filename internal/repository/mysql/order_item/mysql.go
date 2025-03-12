package orderitemrepo

import "gorm.io/gorm"

type mysqlOrderItem struct {
	db *gorm.DB
}

func NewMySQLOrder(db *gorm.DB) *mysqlOrderItem {
	return &mysqlOrderItem{db}
}
