package imagestorage

import "gorm.io/gorm"

var (
	SelectFields = []string{
		"`image`.`image_id`",
		"`image`.`url`",
		"`image`.`alt_text`",
		"`image`.`resource_id`",
		"`image`.`resource_type`",
	}
)

type mysqlImage struct {
	db *gorm.DB
}

func NewMySQLImage(db *gorm.DB) *mysqlImage {
	return &mysqlImage{db}
}
