package categorystorage

import "gorm.io/gorm"

type mysqlCategory struct {
	db *gorm.DB
}

func NewMySQLCategory(db *gorm.DB) *mysqlCategory {
	return &mysqlCategory{db}
}

var (
	SelectFields = []string{
		"category_id",
		"name",
		"description",
		"status",
	}

	AllowedSortFields = map[string]bool{
		"name":        true,
		"created_at":  true,
		"updated_at":  true,
		"category_id": true,
	}

	ImageSelectFields = []string{
		"`image`.`image_id`",
		"`image`.`url`",
		"`image`.`alt_text`",
		"`image`.`resource_id`",
		"`image`.`resource_type`",
	}
)
