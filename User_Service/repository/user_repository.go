package repository

import (
	"User_Service/entity"
	"User_Service/entity/httpentity"
	"User_Service/lib"
	"context"
	"errors"
	"github.com/uptrace/bun"
)

type UserRepository struct {
	DB *bun.DB
}

func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}
func (repo *UserRepository) CreateUser(ctx context.Context, data httpentity.CreateUserRegistration) error {
	//hashing the password

	passwordToHash := data.Password
	HashPassword := lib.HashPassword(passwordToHash)
	data.Password = HashPassword

	_, err := repo.DB.NewInsert().Model(&data).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) GetAllUser(ctx context.Context, filter entity.UserFilter) ([]entity.UserRegistration, int, error) {
	var data []entity.UserRegistration

	query := repo.DB.NewSelect().Model(&data)

	if filter.Keyword != "" {
		query.Where("title Ilike ?", "%"+filter.Keyword+"%")
	}
	if filter.FirstName != "" {
		query.Where("title Ilike ?", "%"+filter.FirstName+"%")
	}
	count, err := query.Limit(filter.GetLimit()).Offset(filter.GetOffset()).ScanAndCount(ctx)
	if err != nil {
		return []entity.UserRegistration{}, 0, err
	}
	return data, count, nil

}
func (repo *UserRepository) GetAUser(ctx context.Context, id int) (entity.UserRegistration, error) {
	var data entity.UserRegistration

	err := repo.DB.NewSelect().Model(&data).Where("user_id =?", id).Scan(ctx)

	if err != nil {
		return entity.UserRegistration{}, err
	}

	return data, nil
}
func (repo *UserRepository) GetAUserRedis(id int) (entity.UserRegistration, error) {
	var data entity.UserRegistration
	ctx := context.Background()
	err := repo.DB.NewSelect().Model(&data).Where("user_id =?", id).Scan(ctx)

	if err != nil {
		return entity.UserRegistration{}, nil
	}

	return data, nil
}

func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (entity.UserRegistration, error) {
	var data entity.UserRegistration

	err := repo.DB.NewSelect().Model(&data).Where("email=?", email).Scan(ctx)
	if err != nil {
		return entity.UserRegistration{}, err
	}
	return data, nil
}
func (repo *UserRepository) UpdateUser(ctx context.Context, data entity.UserRegistration, id int) error {
	_, err := repo.DB.NewUpdate().Model(&data).
		ExcludeColumn("created_at").
		ExcludeColumn("created_by").
		ExcludeColumn("deleted_at").
		ExcludeColumn("updated_by").
		Set("update_at=NOW()").
		Set("first_name=?", data.FirstName).
		Set("last_name=?", data.LastName).
		Set("occupation=?", data.Occupation).
		Set("email=?", data.Email).
		Set("password=?", data.Password).
		Where("user_id=?", id).
		Exec(ctx)

	if err != nil {
		return err
	}
	return nil
}
func (repo *UserRepository) DeleteUser(ctx context.Context, id int) error {
	var data entity.UserRegistration
	res, err := repo.DB.NewDelete().Model(&data).Where("user_id=?", id).Exec(ctx)
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no data matched")
	}
	return err

}
