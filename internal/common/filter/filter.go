package filter

import (
	"client/internal/common/datatypes"
)

type Filter struct {
	Status       datatypes.Status      `json:"status,omitempty" form:"status"`
	Search       string                `json:"search,omitempty" form:"search"`
	MinPrice     float64               `json:"min_price,omitempty" form:"min_price"`
	MaxPrice     float64               `json:"max_price,omitempty" form:"max_price"`
	CategoryID   uint64                `json:"category_id,omitempty" form:"category_id"`
	IDs          []uint64              `json:"ids,omitempty"`
	StartDate    *datatypes.CustomDate `json:"start_date,omitempty" form:"start_date"`
	EndDate      *datatypes.CustomDate `json:"end_date,omitempty" form:"end_date"`
	IngredientID *uint64               `json:"ingredient,omitempty" form:"ingredient_id"`
	ProductID    uint64                `json:"product_id,omitempty"`
	OrderDate
	Product
	ReportOrderDate
}

type OrderDate struct {
	ExpirationDate      *datatypes.CustomDate `json:"expiration_date,omitempty" form:"expiration_date"`
	ReceivedDate        *datatypes.CustomDate `json:"received_date,omitempty" form:"received_date"`
	StartExpirationDate *datatypes.CustomDate `json:"start_expiration_date,omitempty" form:"start_expiration_date"`
	EndExpirationDate   *datatypes.CustomDate `json:"end_expiration_date,omitempty" form:"end_expiration_date"`
	StartReceivedDate   *datatypes.CustomDate `json:"start_received_date,omitempty" form:"start_received_date"`
	EndReceivedDate     *datatypes.CustomDate `json:"end_received_date,omitempty" form:"end_received_date"`
	InOrderDate         *datatypes.CustomDate `json:"order_date,omitempty" form:"order_date"`
}

type ReportOrderDate struct {
	InDate *datatypes.CustomDate `json:"in_date,omitempty" form:"in_date"`
}

type Product struct {
	//ProductIDs []uint64 `json:"product_ids,omitempty"`
	Name string `json:"product_name,omitempty"`
}
