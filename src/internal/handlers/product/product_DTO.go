package product

import "time"

type productCreateRequest struct {
	Id             string    `json:"id"`
	Name           string    `json:"name" validate:"required"`
	Category       string    `json:"category" validate:"required"`
	Price          float64   `json:"price" validate:"required"`
	AvailableStock int       `json:"available_stock" validate:"required"`
	LastUpdate     time.Time `json:"last_update_date"`
	SupplierId     string    `json:"suppler_id" validate:"required"`
	ImageId        string    `json:"image_id"`
	Image          []byte    `json:"image" validate:"required"`
}

//type productReduceRequest struct {
//	Id string `json:"id"`
//
//}
