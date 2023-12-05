//go:generate easyjson -omit_empty=false $GOFILE
package httpentity

//easyjson:json
type Response struct {
	Success bool         `json:"success"`
	Message string       `json:"message,omitempty"`
	Errors  []FieldError `json:"errors,omitempty"`
	Data    interface{}  `json:"data,omitempty"`
}

//easyjson:json
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

//easyjson:json
type ErrorResponse struct {
	Success      bool         `json:"success"`
	ErrorCode    string       `json:"error_code"`
	ErrorMessage string       `json:"error_message,omitempty"`
	FieldErrors  []FieldError `json:"field_errors,omitempty"`
}

//easyjson:json
type OkResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
