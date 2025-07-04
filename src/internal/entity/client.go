package entity

import "time"

type Client struct {
	Id               string    `json:"id"`
	ClientName       string    `json:"client_name"`
	ClientSurname    string    `json:"client_sure_name"`
	BirthDate        time.Time `json:"birth_date"`
	Gender           string    `json:"gender"`
	RegistrationDate time.Time `json:"register_date"`
	AddressId        string    `json:"address_id"`
	//Address   `json:"address,omitempty"`
}
