package dto

import (
	"backend2/internal/entity"
	"time"
)

type ClientCreateRequestDTO struct {
	ClientName    string           `json:"client_name" validate:"required" example:"Harry"`
	ClientSurname string           `json:"client_sure_name" validate:"required" example:"Potter"`
	BirthDate     time.Time        `json:"birth_date" validate:"required" example:"2000-07-31T00:00:00Z"`
	Gender        string           `json:"gender" validate:"required" example:"male"`
	Address       AddressCreateDTO `json:"address"`
}

type ClientUpdateRequestDTO struct {
	Country string `json:"country" validate:"required" example:"UK"`
	City    string `json:"city" validate:"required" example:"London"`
	Street  string `json:"street" validate:"required" example:"Grimmauld Place"`
}

type ClientResponseDTO struct {
	Id               string     `json:"id" example:"f19a3a7-12f5-4332-9582-624519c3eaea"`
	ClientName       string     `json:"client_name" example:"Harry"`
	ClientSurname    string     `json:"client_sure_name" example:"Potter"`
	BirthDate        time.Time  `json:"birth_date" example:"2000-07-31T00:00:00Z"`
	Gender           string     `json:"gender" example:"male"`
	RegistrationDate time.Time  `json:"register_date" example:"2020-09-01T12:00:00Z"`
	AddressId        string     `json:"address_id" example:"a123b456-c789-d012-e345-67890abcdef1"`
	Address          AddressDTO `json:"address"`
}

type ClientsResponseDTO struct {
	Clients []ClientResponseDTO `json:"clients"` //swaggertype:"array,clients"
}

type ClientsNotFound struct {
	Clients []entity.Client `json:"clients"  swaggertype:"array,object"`
	Message string          `json:"message"`
}
