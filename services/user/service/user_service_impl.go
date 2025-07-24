package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/user/model"
	"github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	addressRepo "github.com/SomeHowMicroservice/shm-be/services/user/repository/address"
	measurementRepo "github.com/SomeHowMicroservice/shm-be/services/user/repository/measurement"
	profileRepo "github.com/SomeHowMicroservice/shm-be/services/user/repository/profile"
	roleRepo "github.com/SomeHowMicroservice/shm-be/services/user/repository/role"
	userRepo "github.com/SomeHowMicroservice/shm-be/services/user/repository/user"
	"github.com/google/uuid"
)

type userServiceImpl struct {
	userRepo        userRepo.UserRepository
	roleRepo        roleRepo.RoleRepository
	profileRepo     profileRepo.ProfileRepository
	measurementRepo measurementRepo.MeasurementRepository
	addressRepo     addressRepo.AddressRepository
}

func NewUserService(userRepo userRepo.UserRepository, roleRepo roleRepo.RoleRepository, profileRepo profileRepo.ProfileRepository, measurementRepo measurementRepo.MeasurementRepository, addressRepo addressRepo.AddressRepository) UserService {
	return &userServiceImpl{
		userRepo,
		roleRepo,
		profileRepo,
		measurementRepo,
		addressRepo,
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
		ID:     uuid.NewString(),
		UserID: user.ID,
	}
	if err = s.profileRepo.Create(ctx, profile); err != nil {
		return nil, fmt.Errorf("tạo hồ sơ người dùng thất bại: %w", err)
	}
	// Gán profile rỗng vào phản hồi
	user.Profile = profile

	// Tạo Measurement trống cho người dùng
	measurement := &model.Measurement{
		ID:     uuid.NewString(),
		UserID: user.ID,
	}
	if err := s.measurementRepo.Create(ctx, measurement); err != nil {
		return nil, fmt.Errorf("tạo bảng size người dùng thất bại: %w", err)
	}

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

func (s *userServiceImpl) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin người dùng thất bại: %w", err)
	}
	if user == nil {
		return nil, customErr.ErrUserNotFound
	}

	return user, nil
}

func (s *userServiceImpl) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
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

func (s *userServiceImpl) UpdateUserProfile(ctx context.Context, req *protobuf.UpdateUserProfileRequest) (*model.User, error) {
	updateData := map[string]interface{}{}
	if req.FirstName != "" {
		updateData["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		updateData["last_name"] = req.LastName
	}
	if req.Gender != "" {
		updateData["gender"] = req.Gender
	}
	if req.Dob != "" {
		parsedDob, err := time.Parse("2006-01-02", req.Dob)
		if err != nil {
			return nil, fmt.Errorf("không thể chuyển đổi định dạng: %w", err)
		}
		updateData["dob"] = parsedDob
	}

	if len(updateData) > 0 {
		if err := s.profileRepo.Update(ctx, req.Id, updateData); err != nil {
			if errors.Is(err, customErr.ErrProfileNotFound) {
				return nil, err
			}
			return nil, fmt.Errorf("cập nhật hồ sơ người dùng thất bại: %w", err)
		}
	}
	// Lấy lại user đã cập nhật
	updatedUser, err := s.userRepo.FindByID(ctx, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin người dùng thất bại: %w", err)
	}
	if updatedUser == nil {
		return nil, customErr.ErrProfileNotFound
	}

	return updatedUser, nil
}

func (s *userServiceImpl) GetMeasurementByUserID(ctx context.Context, userID string) (*model.Measurement, error) {
	measurement, err := s.measurementRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("lấy địa chỉ người dùng thất bại: %w", err)
	}
	if measurement == nil {
		return nil, customErr.ErrAddressesNotFound
	}

	return measurement, nil
}

func (s *userServiceImpl) UpdateUserMeasurement(ctx context.Context, req *protobuf.UpdateUserMeasurementRequest) (*model.Measurement, error) {
	measurement, err := s.measurementRepo.FindByID(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("lấy độ đo người dùng thất bại: %w", err)
	}
	if measurement == nil {
		return nil, customErr.ErrMeasurementNotFound
	}

	if measurement.UserID != req.UserId {
		return nil, customErr.ErrForbidden
	}

	updateData := map[string]interface{}{}
	if req.Height != 0 {
		updateData["height"] = req.Height
	}
	if req.Weight != 0 {
		updateData["weight"] = req.Weight
	}
	if req.Chest != 0 {
		updateData["chest"] = req.Chest
	}
	if req.Waist != 0 {
		updateData["waist"] = req.Waist
	}
	if req.Butt != 0 {
		updateData["butt"] = req.Butt
	}
	if len(updateData) >= 0 {
		if err := s.measurementRepo.Update(ctx, measurement.ID, updateData); err != nil {
			if errors.Is(err, customErr.ErrMeasurementNotFound) {
				return nil, err
			}
			return nil, fmt.Errorf("cập nhật độ đo người dùng thất bại: %w", err)
		}

		measurement, err = s.measurementRepo.FindByID(ctx, measurement.ID)
		if err != nil {
			return nil, fmt.Errorf("lấy độ đo người dùng thất bại: %w", err)
		}
		if measurement == nil {
			return nil, customErr.ErrMeasurementNotFound
		}
	}

	return measurement, nil
}

func (s *userServiceImpl) GetAddressesByUserID(ctx context.Context, userID string) ([]*model.Address, error) {
	addresses, err := s.addressRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("lấy địa chỉ người dùng thất bại: %w", err)
	}

	return addresses, nil
}

func (s *userServiceImpl) CreateAddress(ctx context.Context, req *protobuf.CreateAddressRequest) (*model.Address, error) {
	countAddreses, err := s.addressRepo.CountByUserID(ctx, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("lấy địa chỉ người dùng thất bại: %w", err)
	}
	if countAddreses == 0 {
		req.IsDefault = true
	}

	address := &model.Address{
		ID:          uuid.NewString(),
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Street:      req.Street,
		Ward:        req.Ward,
		Province:    req.Province,
		IsDefault:   req.IsDefault,
		UserID:      req.UserId,
	}
	if err = s.addressRepo.Create(ctx, address); err != nil {
		return nil, fmt.Errorf("tạo địa chỉ người dùng thất bại: %w", err)
	}

	if address.IsDefault {
		addresses, err := s.addressRepo.FindByUserID(ctx, req.UserId)
		if err != nil {
			return nil, fmt.Errorf("lấy địa chỉ người dùng thất bại: %w", err)
		}

		for _, addr := range addresses {
			if addr.ID != address.ID && addr.IsDefault {
				if err := s.addressRepo.Update(ctx, addr.ID, map[string]interface{}{"is_default": false}); err != nil {
					return nil, fmt.Errorf("cập nhật địa chỉ mặc định thất bại: %w", err)
				}
			}
		}
	}

	return address, nil
}

func (s *userServiceImpl) UpdateAddress(ctx context.Context, req *protobuf.UpdateAddressRequest) (*model.Address, error) {
	address, err := s.addressRepo.FindByID(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("tìm địa chỉ thất bại: %w", err)
	}
	if address == nil {
		return nil, customErr.ErrAddressesNotFound
	}

	if address.UserID != req.UserId {
		return nil, customErr.ErrForbidden
	}

	countAddreses, err := s.addressRepo.CountByUserID(ctx, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("đếm địa chỉ thất bại: %w", err)
	}

	updateData := map[string]interface{}{}
	if req.FullName != address.FullName {
		updateData["full_name"] = req.FullName
	}
	if req.PhoneNumber != address.PhoneNumber {
		updateData["phone_number"] = req.PhoneNumber
	}
	if req.Ward != address.Ward {
		updateData["ward"] = req.Ward
	}
	if req.Street != address.Street {
		updateData["street"] = req.Street
	}
	if req.Province != address.Province {
		updateData["province"] = req.Province
	}

	if address.IsDefault != req.IsDefault {
		if !req.IsDefault {
			if countAddreses == 1 {
				req.IsDefault = true
			} else if countAddreses > 1 {
				addresses, err := s.addressRepo.FindByUserID(ctx, req.UserId)
				if err != nil {
					return nil, fmt.Errorf("tìm địa chỉ người dùng thất bại: %w", err)
				}

				if err = s.addressRepo.Update(ctx, addresses[1].ID, map[string]interface{}{"is_default": true}); err != nil {
					return nil, fmt.Errorf("không thể cập nhật địa chỉ %w", err)
				}
			}
		} else {
			addresses, err := s.addressRepo.FindByUserID(ctx, req.UserId)
			if err != nil {
				return nil, fmt.Errorf("tìm địa chỉ người dùng thất bại: %w", err)
			}

			for _, addr := range addresses {
				if addr.ID != address.ID {
					if err := s.addressRepo.Update(ctx, addr.ID, map[string]interface{}{"is_default": false}); err != nil {
						return nil, fmt.Errorf("cập nhật địa chỉ mặc định thất bại: %w", err)
					}
				}
			}
		}
		updateData["is_default"] = req.IsDefault
	}

	if len(updateData) > 0 {
		if err := s.addressRepo.Update(ctx, address.ID, updateData); err != nil {
			return nil, fmt.Errorf("lỗi cập nhật địa chỉ: %w", err)
		}

		address, err = s.addressRepo.FindByID(ctx, req.Id)
		if err != nil {
			return nil, fmt.Errorf("không thể tìm địa chỉ: %w", err)
		}
		if address == nil {
			return nil, customErr.ErrAddressesNotFound
		}
	}

	return address, nil
}

func (s *userServiceImpl) DeleteAddress(ctx context.Context, req *protobuf.DeleteAddressRequest) error {
	address, err := s.addressRepo.FindByID(ctx, req.Id)
	if err != nil {
		return fmt.Errorf("tìm địa chỉ thất bại: %w", err)
	}
	if address == nil {
		return customErr.ErrAddressesNotFound
	}

	if address.UserID != req.UserId {
		return customErr.ErrForbidden
	}

	if err := s.addressRepo.Delete(ctx, req.Id); err != nil {
		if errors.Is(err, customErr.ErrAddressesNotFound) {
			return err
		}
		return fmt.Errorf("xóa địa chỉ thất bại: %w", err)
	}

	if address.IsDefault {
		addresses, err := s.addressRepo.FindByUserID(ctx, req.UserId)
		if err != nil {
			return fmt.Errorf("tìm địa chỉ người dùng thất bại: %w", err)
		}

		if len(addresses) >= 1 {
			if err := s.addressRepo.Update(ctx, addresses[0].ID, map[string]interface{}{"is_default": true}); err != nil {
				return fmt.Errorf("cập nhật địa chỉ mặc định thất bại: %w", err)
			}
		}
	}

	return nil
}
