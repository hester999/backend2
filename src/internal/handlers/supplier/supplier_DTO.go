package supplier

type supplierCreateRequest struct {
	//Id          string `json:"id"`
	Name string `json:"name" validate:"required"`
	//AddressId   string `json:"address"`
	PhoneNumber string `json:"phone" validate:"required"`
	City        string `json:"city" validate:"required"`
	Street      string `json:"street" validate:"required"`
	Country     string `json:"country" validate:"required"`
}

type supplierUpdateAddressRequest struct {
	//AddressId string `json:"address_id" validate:"required"`
	City    string `json:"city" validate:"required"`
	Street  string `json:"street" validate:"required"`
	Country string `json:"country" validate:"required"`
}
