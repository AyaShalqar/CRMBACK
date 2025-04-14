package shop

import "time"

type Item struct {
	ID            int       `json:"id" db:"id"`
	ShopID        int       `json:"shop_id" db:"shop_id"`
	Name          string    `json:"name" db:"name"`
	Brand         string    `json:"brand" db:"brand"`
	Category      string    `json:"category" db:"category"`
	Size          string    `json:"size" db:"size"`
	PurchasePrice float64   `json:"purchase_price" db:"purchase_price"`
	SalePrice     float64   `json:"sale_price" db:"sale_price"`
	PhotoURL      string    `json:"photo_url" db:"photo_url"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type CreateItemRequest struct {
	Name          string  `json:"name"`
	Brand         string  `json:"brand"`
	Category      string  `json:"category"`
	Size          string  `json:"size"`
	PurchasePrice float64 `json:"purchase_price"`
	SalePrice     float64 `json:"sale_price"`
	PhotoURL      string  `json:"photo_url"`
}
