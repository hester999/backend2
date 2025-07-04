package dto

type AddressDTO struct {
	ID      string `json:"id" example:"a123b456-c789-d012-e345-67890abcdef1"`
	Country string `json:"country" example:"UK" validate:"required"`
	City    string `json:"city" example:"London" validate:"required"`
	Street  string `json:"street" example:"Privet Drive" validate:"required"`
}

type AddressCreateDTO struct {
	Country string `json:"country" example:"UK" validate:"required"`
	City    string `json:"city" example:"London" validate:"required"`
	Street  string `json:"street" example:"Privet Drive" validate:"required"`
}
