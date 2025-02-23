package imagemodel

var (
	EntityName = "image"
)

type Image struct {
	ImageID uint64 `gorm:"column:image_id;primaryKey;autoIncrement" json:"image_id"`
	URL     string `gorm:"column:url;size:300;not null" json:"url"`
	AltText string `gorm:"column:alt_text;size:255" json:"alt_text,omitempty"`

	ResourceID   uint64 `gorm:"column:resource_id;index" json:"resource_id"`
	ResourceType string `gorm:"column:resource_type;size:50;index" json:"resource_type"`
	//
	//	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at,omitempty"`
	//	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at,omitempty"`
}

func (Image) TableName() string {
	return "image"
}
