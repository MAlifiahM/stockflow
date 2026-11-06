package domain

type Error struct {
	Code    int      `json:"code"`
	Errors  []string `json:"errors,omitempty"`
	Message string   `json:"message"`
}

type ErrorDetail struct {
	Code    string           `json:"code,omitempty"`
	Message string           `json:"message,omitempty"`
	Error   *ErrorCodeDetail `json:"error,omitempty"`
}

type ErrorCodeDetail struct {
	Code      string  `json:"code"`
	SKU       *string `json:"sku,omitempty"`
	Available *int32  `json:"available,omitempty"`
}

func NewError(code int, message string) Error {
	return Error{Code: code, Message: message}
}

func (e Error) Error() string {
	return e.Message
}
