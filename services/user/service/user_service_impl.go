package service

import (
	"context"
	"time"

	"github.com/SomeHowMicroservice/shm-be/services/user/model"
	"github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"github.com/SomeHowMicroservice/shm-be/services/user/repository"
	"github.com/SomeHowMicroservice/shm-be/services/user/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userServiceImpl struct {
	repo repository.UserRepository
	protobuf.UnimplementedUserServiceServer
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{
		repo: repo,
		UnimplementedUserServiceServer: protobuf.UnimplementedUserServiceServer{},
	}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, req *protobuf.CreateUserRequest) (*protobuf.UserResponse, error) {
	exists, err := s.repo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Không thể kiểm tra username: %v", err)
	}
	if exists {
		return nil, status.Error(codes.AlreadyExists, "Username đã tồn tại")
	}

	exists, err = s.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Không thể kiểm tra email: %v", err)
	}
	if exists {
		return nil, status.Error(codes.AlreadyExists, "Email đã tồn tại")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	user := &model.User{
		Username: req.Username,
		Email: req.Email,
		Password: hashedPassword,
	}

	if err = s.repo.Create(ctx, user); err != nil {
		return nil, status.Errorf(codes.Internal, "Không thể tạo người dùng: %v", err)
	}

	return &protobuf.UserResponse{
		Id: user.ID.String(),
		Username: user.Username,
		Email: user.Email,
		Password: hashedPassword,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}
