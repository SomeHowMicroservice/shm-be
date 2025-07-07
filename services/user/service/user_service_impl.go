package service

import (
	"context"
	"fmt"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/user/model"
	"github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"github.com/SomeHowMicroservice/shm-be/services/user/repository"
	"github.com/SomeHowMicroservice/shm-be/services/user/utils"
)

type userServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{
		repo: repo,
	}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, req *protobuf.CreateUserRequest) (*model.User, error) {
	exists, err := s.repo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, customErr.ErrUsernameAlreadyExists
	}

	exists, err = s.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, customErr.ErrEmailAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Email: req.Email,
		Password: hashedPassword,
	}

	if err = s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("không thể tạo người dùng: %w", err)
	}

	return user, nil
}
