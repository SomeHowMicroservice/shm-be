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
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

func NewUserService(userRepo repository.UserRepository, roleRepo repository.RoleRepository) UserService {
	return &userServiceImpl{
		userRepo,
		roleRepo,
	}
}

func (s *userServiceImpl) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return false, fmt.Errorf("kiểm tra email thất bại: %w", err)
	}
	return exists, nil
}

func (s *userServiceImpl) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	exists, err := s.userRepo.ExistsByUsername(ctx, username)
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
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("tạo người dùng thất bại: %w", err)
	}
	// Lấy quyền user để xét
	role, err := s.roleRepo.FindByName(ctx, model.RoleUser)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin quyền thất bại: %w", err)
	}
	if role == nil {
		return nil, customErr.ErrRoleNotFound
	}
	// Thêm quyền cho người dùng
	if err = s.roleRepo.CreateUserRoles(ctx, user.ID, role.ID); err != nil {
		return nil, fmt.Errorf("thêm quyền cho người dùng thất bại: %w", err)
	}
	// Gán quyền vào phản hổi
	user.Roles = []model.Role{*role}
	return user, nil
}

func (s *userServiceImpl) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin người dùng thất bại: %w", err)
	}
	if user == nil {
		return nil, customErr.ErrUserNotFound
	}
	return user, nil
}
