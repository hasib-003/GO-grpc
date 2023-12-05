package entity

type Response struct {
	Success bool         `json:"success"`
	Message string       `json:"message,omitempty"`
	Errors  []FieldError `json:"errors,omitempty"`
	Data    interface{}  `json:"data,omitempty"`
}

type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type RegistrationResponse struct {
	OtpRequired int  `json:"otp_required"`
	OtpSent     bool `json:"otp_sent"`
}
