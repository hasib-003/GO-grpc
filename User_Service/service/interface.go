package service

import (
	"User_Service/entity"
	"User_Service/entity/httpentity"
	"context"
)

type UserRepository interface {
	CreateUser(ctx context.Context, data httpentity.CreateUserRegistration) error
	GetAllUser(ctx context.Context, filter entity.UserFilter) ([]entity.UserRegistration, int, error)
	GetAUser(ctx context.Context, id int) (entity.UserRegistration, error)
	UpdateUser(ctx context.Context, data entity.UserRegistration, id int) error
	DeleteUser(ctx context.Context, id int) error
	GetUserByEmail(ctx context.Context, email string) (entity.UserRegistration, error)
}
