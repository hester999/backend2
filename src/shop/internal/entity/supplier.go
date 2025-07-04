package entity

//{
//id
//name
//address_id
//phone_number
//}

type Supplier struct {
	Id          string
	Name        string
	AddressId   string
	PhoneNumber string
	Address     Address
}
