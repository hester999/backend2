package client

import "time"

type clientCreateRequest struct {
	Id               string    `json:"id" `
	ClientName       string    `json:"client_name" validate:"required"`
	ClientSurname    string    `json:"client_sure_name" validate:"required"`
	BirthDate        time.Time `json:"birth_date" validate:"required"`
	Gender           string    `json:"gender" validate:"required"`
	RegistrationDate time.Time `json:"register_date"`
	AddressId        string    `json:"address_id"`
	Address          struct {
		ID      string `json:"id"`
		Country string `json:"country" validate:"required"`
		City    string `json:"city" validate:"required"`
		Street  string `json:"street" validate:"required"`
	}
}

type clientUpdateRequest struct {
	//Id      string `json:"id" validate:"required"`
	Country string `json:"country" validate:"required"`
	City    string `json:"city" validate:"required"`
	Street  string `json:"street" validate:"required"`
}

type clientGetByNameSurnameRequest struct {
	Name    string `json:"name" validate:"required"`
	Surname string `json:"surname" validate:"required"`
}

type clientRequestById struct {
	Id string `json:"id"`
}

type responseErr struct {
	Status  int
	Message string
}
