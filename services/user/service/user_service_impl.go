package service

import (
	"context"
	"fmt"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/user/model"
	"github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"github.com/SomeHowMicroservice/shm-be/services/user/repository"
	"github.com/google/uuid"
)

type userServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{repo}
}

func (s *userServiceImpl) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	exists, err := s.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return false, fmt.Errorf("kiểm tra email thất bại: %w", err)
	}
	return exists, nil
}

func (s *userServiceImpl) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	exists, err := s.repo.ExistsByUsername(ctx, username)
	if err != nil {
		return false, fmt.Errorf("kiểm tra username thất bại: %w", err)
	}
	return exists, nil
}

func (s *userServiceImpl) CreateUser(ctx context.Context, req *protobuf.CreateUserRequest) (*model.User, error) {
	user := &model.User{
		ID:       uuid.NewString(),
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("tạo người dùng thất bại: %w", err)
	}
	return user, nil
}

func (s *userServiceImpl) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if user == nil {
		return nil, customErr.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin người dùng thất bại: %w", err)
	}
	return user, nil
}
