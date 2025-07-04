package entity

import "time"

type Client struct {
	Id               string
	ClientName       string
	ClientSurname    string
	BirthDate        time.Time
	Gender           string
	RegistrationDate time.Time
	AddressId        string
	Address
}
