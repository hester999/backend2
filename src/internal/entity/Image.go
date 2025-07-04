package entity

//}images
//{
//id : UUID
//image: bytea
//}

type Image struct {
	Id    string `json:"id"`
	Image []byte `json:"image"`
}
