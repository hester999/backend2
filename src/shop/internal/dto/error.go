package dto

type ErrorResponse struct {
	Message string `json:"status" example:"error"`
	Code    int    `json:"code"`
}

type Error400 struct {
	Message string `json:"status" example:"invalid JSON"`
	Code    int    `json:"code" example:"400"`
}

type Error404 struct {
	Message string `json:"status" example:"not found"`
	Code    int    `json:"code" example:"404"`
}

type Error500 struct {
	Message string `json:"status" example:"internal server error"`
	Code    int    `json:"code" example:"500"`
}
