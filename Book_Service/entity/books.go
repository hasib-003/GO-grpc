package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type Books struct {
	bun.BaseModel   `bun:"table:books"`
	Id              int        `json:"id" bun:",pk,autoincrement"`
	Title           string     `json:"title" bun:",notnull"`
	UserId          int        `json:"user_id" bun:",notnull"`
	PublicationYear int        `json:"publication_year" bun:",notnull"`
	CreatedAt       time.Time  `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	UpdateAt        *time.Time `json:"update_at" bun:",nullzero"`
	DeletedAt       *time.Time `json:"-" bun:",soft_delete"`
	CreatedBy       *string    `json:"created_by" bun:"type:uuid,default:uuid_generate_v4()"`
	UpdatedBy       *string    `json:"updated_by" bun:"type:uuid,default:uuid_generate_v4()"`
}

func (p *Books) Validate() []FieldError {
	return validate(p)
}

type ListAllBooksResponse struct {
	Total int     `json:"total"`
	Page  int     `json:"page"`
	Books []Books `json:"books"`
}

type BooksFilter struct {
	Keyword         string `query:"keyword" validate:"omitempty"`
	Author          string `query:"author"  validate:"omitempty"`
	PublicationYear int    `query:"publication_year" validate:"omitempty"`
	PaginationRequest
}
