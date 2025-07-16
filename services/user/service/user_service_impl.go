package service

import (
	"context"
	"errors"
	"fmt"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/user/model"
	"github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	profileRepo "github.com/SomeHowMicroservice/shm-be/services/user/repository/profile"
	roleRepo "github.com/SomeHowMicroservice/shm-be/services/user/repository/role"
	userRepo "github.com/SomeHowMicroservice/shm-be/services/user/repository/user"
	"github.com/google/uuid"
)

type userServiceImpl struct {
	userRepo    userRepo.UserRepository
	roleRepo    roleRepo.RoleRepository
	profileRepo profileRepo.ProfileRepository
}

func NewUserService(userRepo userRepo.UserRepository, roleRepo roleRepo.RoleRepository, profileRepo profileRepo.ProfileRepository) UserService {
	return &userServiceImpl{
		userRepo,
		roleRepo,
		profileRepo,
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
	user.Roles = []*model.Role{role}
	// Tạo profile trống cho người dùng
	profile := &model.Profile{
		ID: uuid.NewString(),
		UserID: user.ID,
	}
	if err = s.profileRepo.Create(ctx, profile); err != nil {
		return nil, fmt.Errorf("tạo hồ sơ người dùng thất bại: %w", err)
	}
	// Gán profile rỗng vào phản hồi
	user.Profile = profile
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

func (s *userServiceImpl) GetUserById(ctx context.Context, id string) (*model.User, error) {
	user, err := s.userRepo.FindById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin người dùng thất bại: %w", err)
	}
	if user == nil {
		return nil, customErr.ErrUserNotFound
	}
	return user, nil
}

func (s *userServiceImpl) UpdateUserPassword(ctx context.Context, req *protobuf.UpdateUserPasswordRequest) error {
	if err := s.userRepo.UpdatePassword(ctx, req.Id, req.NewPassword); err != nil {
		if errors.Is(err, customErr.ErrUserNotFound) {
			return err
		}
		return fmt.Errorf("cập nhật mật khẩu thất bại: %w", err)
	}
	return nil
}
