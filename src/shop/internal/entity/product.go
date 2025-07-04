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
	Id             string
	Name           string
	Category       string
	Price          float64
	AvailableStock int
	LastUpdate     time.Time
	SupplierId     string
	ImageId        string
}
