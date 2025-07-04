package dto

import "backend2/internal/entity"

type SupplierCreateRequestDTO struct {
	Name        string           `json:"name" validate:"required" example:"Magic Supplies Inc."`
	PhoneNumber string           `json:"phone" validate:"required" example:"+44-123-456-789"`
	Address     AddressCreateDTO `json:"address"`
}

type SupplierUpdateAddressRequestDTO struct {
	City    string `json:"city" validate:"required" example:"Edinburgh"`
	Street  string `json:"street" validate:"required" example:"Royal Mile"`
	Country string `json:"country" validate:"required" example:"UK"`
}

type SupplierResponseDTO struct {
	Id          string     `json:"id" example:"supplier-1234"`
	Name        string     `json:"name" example:"Magic Supplies Inc."`
	PhoneNumber string     `json:"phone" example:"+44-123-456-789"`
	AddressId   string     `json:"address_id" example:"address-5678"`
	Address     AddressDTO `json:"address"`
}

type SuppliersResponse struct {
	Suppliers []SupplierResponseDTO `json:"suppliers"`
}

type SuppliersNotFound struct {
	Suppliers []entity.Supplier `json:"suppliers" swaggertype:"array,object"`
	Message   string            `json:"message"`
}
