package service

import (
	"Book_Service/entity"
	"context"
)

type Service struct {
	BooksRepository Repository
}

func NewBooksService(booksRepository Repository) *Service {
	return &Service{
		BooksRepository: booksRepository,
	}
}

func (s *Service) Create(ctx context.Context, data entity.Books) error {

	err := s.BooksRepository.Create(ctx, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *Service) ListAllBooks(ctx context.Context, filter entity.BooksFilter) ([]entity.Books, int, error) {
	res, count, err := s.BooksRepository.ListAllBooks(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return res, count, nil
}

func (s *Service) GetABook(ctx context.Context, bookId string) (entity.Books, error) {
	res, err := s.BooksRepository.GetABook(ctx, bookId)
	if err != nil {
		return entity.Books{}, err
	}
	return res, nil
}

func (s *Service) Update(ctx context.Context, data entity.Books, bookId string) error {
	err := s.BooksRepository.Update(ctx, data, bookId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, bookId string) error {
	err := s.BooksRepository.Delete(ctx, bookId)
	if err != nil {
		return err
	}
	return nil
}

//func (s *Service) GetBooksByUserID(ctx context.Context, request entity.MessageRequest) (entity.MessageResponse, error) {
//	books, err := s.BooksRepository.GetBooksByUserID(ctx, request.UserID)
//	if err != nil {
//		return entity.MessageResponse{}, err
//	}
//	response := entity.MessageResponse{
//		Books: books,
//	}
//
//	return response, nil
//}
