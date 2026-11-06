package domain

type Response struct {
	Data interface{} `json:"data,omitempty"`
}

type PaginationResponse struct {
	Data       interface{} `json:"data,omitempty"`
	Page       int64       `json:"page"`
	Limit      int64       `json:"limit"`
	TotalPages int64       `json:"total_pages"`
	TotalItems int64       `json:"total_items"`
}
