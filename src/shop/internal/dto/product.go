package dto

import (
	"backend2/internal/entity"
	"time"
)

type ProductCreateRequest struct {
	Name           string  `json:"name" validate:"required" example:"Potion of Healing"`
	Category       string  `json:"category" validate:"required" example:"Alchemy"`
	Price          float64 `json:"price" validate:"required" example:"49.99"`
	AvailableStock int     `json:"available_stock" validate:"required" example:"120"`
	SupplierId     string  `json:"suppler_id" validate:"required" example:"supplier-abc-123"`
}

type ProductResponse struct {
	Id             string    `json:"id" example:"product-xyz-789"`
	Name           string    `json:"name" example:"Potion of Healing"`
	Category       string    `json:"category" example:"Alchemy"`
	Price          float64   `json:"price" example:"49.99"`
	AvailableStock int       `json:"available_stock" example:"120"`
	LastUpdate     time.Time `json:"last_update_date" example:"2025-07-01T15:04:05Z"`
	SupplierId     string    `json:"suppler_id" example:"supplier-abc-123"`
	ImageId        string    `json:"image_id" example:"img-00112233"`
}

type ProductsResponse struct {
	Products []ProductResponse `json:"products"`
}

type ProductsNotFound struct {
	Products []entity.Product `json:"products" swaggertype:"array,object"`
	Message  string           `json:"message"`
}
