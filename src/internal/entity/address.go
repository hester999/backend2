package entity

//{
//id
//country
//city
//street
//}

type Address struct {
	ID      string `json:"id"`
	Country string `json:"country"`
	City    string `json:"city"`
	Street  string `json:"street"`
}
