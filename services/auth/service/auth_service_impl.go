package service

import (
	"context"
	"fmt"
	"time"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/auth/common"
	"github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
	"github.com/SomeHowMicroservice/shm-be/services/auth/repository"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"github.com/google/uuid"
)

type authServiceImpl struct {
	cacheRepo  repository.CacheRepository
	userClient userpb.UserServiceClient
}

func NewAuthService(cacheRepo repository.CacheRepository, userClient userpb.UserServiceClient) AuthService {
	return &authServiceImpl{
		userClient: userClient,
		cacheRepo:  cacheRepo,
	}
}

func (s *authServiceImpl) SignUp(ctx context.Context, req *protobuf.SignUpRequest) (string, error) {
	uRes, err := s.userClient.CheckUsernameExists(ctx, &userpb.CheckUsernameExistsRequest{Username: req.Username})
	if err != nil {
		return "", err
	}
	if uRes.Exists {
		return "", customErr.ErrUsernameAlreadyExists
	}

	eRes, err := s.userClient.CheckEmailExists(ctx, &userpb.CheckEmailExistsRequest{Email: req.Email})
	if err != nil {
		return "", err
	}
	if eRes.Exists {
		return "", customErr.ErrEmailAlreadyExists
	}

	otp := common.GenerateOTP(6)
	registrationToken := uuid.NewString()

	hashedPassword, err := common.HashPassword(req.Password)
	if err != nil {
		return "", fmt.Errorf("băm mật khẩu thất bại: %w", err)
	}

	registrationData := common.RegistrationData{
		Email:    req.Email,
		Username: req.Username,
		Password: hashedPassword,
		Otp:      otp,
		Attempts: 0,
	}

	if err = s.cacheRepo.SaveRegistrationData(ctx, registrationToken, registrationData, 3*time.Minute); err != nil {
		return "", fmt.Errorf("lưu dữ liệu đăng ký vào bộ nhớ đệm thất bại: %w", err)
	}

	return registrationToken, nil
}
