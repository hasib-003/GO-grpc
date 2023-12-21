package repository

import (
	"Book_Service/entity"
	"context"
	"errors"
	"github.com/uptrace/bun"
)

type BooksRepository struct {
	DB *bun.DB
}

func NewBooksRepository(db *bun.DB) *BooksRepository {
	return &BooksRepository{
		DB: db,
	}
}

func (repo *BooksRepository) Create(ctx context.Context, data entity.Books) error {
	_, err := repo.DB.NewInsert().Model(&data).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *BooksRepository) ListAllBooks(ctx context.Context, filter entity.BooksFilter) ([]entity.Books, int, error) {
	var data []entity.Books

	query := repo.DB.NewSelect().Model(&data)

	if filter.Keyword != "" {
		query.Where("title ILike ?", "%"+filter.Keyword+"%")
	}

	if filter.Author != "" {
		query.Where("author ILike ?", "%"+filter.Author+"%")
	}

	if filter.PublicationYear != 0 {
		query.Where("publication_year =?", filter.PublicationYear)
	}

	count, err := query.Limit(filter.GetLimit()).Offset(filter.GetOffset()).ScanAndCount(ctx)

	if err != nil {
		return []entity.Books{}, 0, err
	}

	return data, count, nil
}

func (repo *BooksRepository) GetABook(ctx context.Context, bookId string) (entity.Books, error) {
	var data entity.Books

	err := repo.DB.NewSelect().Model(&data).Where("id =?", bookId).Scan(ctx)

	if err != nil {
		return entity.Books{}, err
	}

	return data, nil
}

func (repo *BooksRepository) Update(ctx context.Context, data entity.Books, bookId string) error {
	_, err := repo.DB.NewUpdate().
		Model(&data).
		ExcludeColumn("created_at").
		ExcludeColumn("created_by").
		ExcludeColumn("deleted_at").
		ExcludeColumn("updated_by").
		Set("update_at = NOW()").
		Set("title = ?", data.Title).
		Set("user_id = ?", data.UserId).
		Set("publication_year = ?", data.PublicationYear).
		Where("id=?", bookId).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *BooksRepository) Delete(ctx context.Context, bookId string) error {
	var data entity.Books

	res, err := repo.DB.NewDelete().Model(&data).Where("id=?", bookId).Exec(ctx)
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no data matched")
	}

	return err
}
func (repo *BooksRepository) GetBooksByUserID(ctx context.Context, userID int) (entity.MessageResponse, error) {
	var books entity.MessageResponse
	err := repo.DB.NewSelect().Model(&books).
		Relation("User").
		Where("user_id = ?", userID).
		Scan(ctx)

	if err != nil {
		return entity.MessageResponse{}, err
	}

	return books, nil
}
