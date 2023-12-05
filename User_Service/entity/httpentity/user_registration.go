package httpentity

import "github.com/uptrace/bun"

type CreateUserRegistration struct {
	bun.BaseModel `bun:"table:user_registration"`

	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name" validate:"required"`
	Occupation string `json:"occupation" validate:"required"`
	Role       string `json:"role" validate:"required"`
	Email      string `json:"email" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

func (p *CreateUserRegistration) Validate() []FieldError {
	return validate(p)
}
