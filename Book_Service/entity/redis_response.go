package entity

import (
	"time"
)

type RedisResponse struct {
	Id              int        `json:"id" `
	Title           string     `json:"title"`
	UserName        string     `json:"user_name" `
	PublicationYear int        `json:"publication_year"`
	CreatedAt       time.Time  `json:"created_at" `
	UpdateAt        *time.Time `json:"update_at" `
	DeletedAt       *time.Time `json:"-"`
	CreatedBy       *string    `json:"created_by" `
	UpdatedBy       *string    `json:"updated_by"`
}
