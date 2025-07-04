package entity

//{
//id
//name
//address_id
//phone_number
//}

type Supplier struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	AddressId   string `json:"address_id"`
	PhoneNumber string `json:"phone"`
}
