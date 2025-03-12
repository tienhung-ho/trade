package orderrepo

import "gorm.io/gorm"

var (
	AllowedSortFields = map[string]bool{
		"order_id":   true,
		"created_at": true,
		"updated_at": true,
		"size":       true,
		"cost":       true,
		"order_date": true,
	}
)

type mysqlOrder struct {
	db *gorm.DB
}

func NewMySQLOrder(db *gorm.DB) *mysqlOrder {
	return &mysqlOrder{db}
}
