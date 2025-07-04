package entity

import "time"

// {
// id
// name
// category
// price
// available_stock // число закупленных экземпляров товара
// last_update_date // число последней закупки
// supplier_id
// image_id: UUID
// }
type Product struct {
	Id             string    `json:"id"`
	Name           string    `json:"name"`
	Category       string    `json:"category"`
	Price          float64   `json:"price"`
	AvailableStock int       `json:"available_stock"`
	LastUpdate     time.Time `json:"last_update_date"`
	SupplierId     string    `json:"suppler_id"`
	ImageId        string    `json:"image_id"`
}
