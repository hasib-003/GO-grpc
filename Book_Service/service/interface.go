package service

import (
	"Book_Service/entity"
	"context"
)

type Repository interface {
	Create(ctx context.Context, data entity.Books) error
	ListAllBooks(ctx context.Context, filter entity.BooksFilter) ([]entity.Books, int, error)
	GetABook(ctx context.Context, bookId string) (entity.Books, error)
	Update(ctx context.Context, data entity.Books, bookId string) error
	Delete(ctx context.Context, bookId string) error
}
