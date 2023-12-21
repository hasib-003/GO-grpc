package service

import (
	"User_Service/entity"
	"User_Service/entity/httpentity"
	"context"
	"fmt"
)

type UserService struct {
	UserRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		UserRepository: userRepository,
	}
}

func (s *UserService) CreateUser(ctx context.Context, data httpentity.CreateUserRegistration) error {
	err := s.UserRepository.CreateUser(ctx, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *UserService) GetAllUser(ctx context.Context, filter entity.UserFilter) ([]entity.UserRegistration, int, error) {
	res, count, err := s.UserRepository.GetAllUser(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return res, count, nil
}

func (s *UserService) GetAUser(ctx context.Context, id int) (entity.UserRegistration, error) {
	res, err := s.UserRepository.GetAUser(ctx, id)
	if err != nil {
		return entity.UserRegistration{}, err
	}
	return res, nil
}
func (s *UserService) GetAUserRedis(id int) (entity.UserRegistration, error) {
	res, err := s.UserRepository.GetAUserRedis(id)
	if err != nil {
		return entity.UserRegistration{}, err
	}
	return res, nil
}
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (entity.UserRegistration, error) {

	fmt.Println(email)

	res, err := s.UserRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return entity.UserRegistration{}, err
	}
	return res, nil
}
func (s *UserService) UpdateUser(ctx context.Context, data entity.UserRegistration, id int) error {
	err := s.UserRepository.UpdateUser(ctx, data, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	err := s.UserRepository.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
