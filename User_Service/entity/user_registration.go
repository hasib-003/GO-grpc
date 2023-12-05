package entity

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/uptrace/bun"
	"time"
)

type UserRegistration struct {
	bun.BaseModel `bun:"table:user_registration"`

	UserId     int        `json:"user_id"bun:",pk,autoincrement"`
	FirstName  string     `json:"first_name"bun:"first_name"`
	LastName   string     `json:"last_name" bun:"last_name"`
	Occupation string     `json:"occupation"bun:"occupation"`
	Role       string     `json:"role"bun:"role"`
	Email      string     `json:"email"bun:"email"`
	Password   string     `json:"password"bun:"password"`
	CreatedAt  time.Time  `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	UpdateAt   *time.Time `json:"update_at" bun:",nullzero"`
	DeletedAt  *time.Time `json:"-" bun:",soft_delete"`
	CreatedBy  *string    `json:"created_by" bun:"type:uuid,default:uuid_generate_v4()"`
	UpdatedBy  *string    `json:"updated_by" bun:"type:uuid,default:uuid_generate_v4()"`
}

func (p *UserRegistration) Validate() []FieldError {
	return validate(p)
}

func (p *UserRegistration) GetJwt(expirationTime time.Time, id int, role string) (*string, error) {
	jwtSecret := "secret"
	fmt.Println("id ", id)
	claims := &JwtClaims{
		UserId: id,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	fmt.Println("claims.............................", claims)
	tokenBase := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenBase.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}
	return &token, nil
}

type UserFilter struct {
	Keyword   string `query:"keyword"`
	FirstName string `query:"first_name"`
	//LastName   string `query:"last_name"`
	//Occupation string `query:"occupation"`
	PaginationRequest
}

type GetAllUserResponses struct {
	Total int                `json:"total"`
	Page  int                `json:"page"`
	Users []UserRegistration `json:"users"`
}
