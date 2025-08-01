package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/common/smtp"
	"github.com/SomeHowMicroservice/shm-be/services/auth/common"
	"github.com/SomeHowMicroservice/shm-be/services/auth/config"
	"github.com/SomeHowMicroservice/shm-be/services/auth/mq"
	"github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
	"github.com/SomeHowMicroservice/shm-be/services/auth/repository"
	"github.com/SomeHowMicroservice/shm-be/services/auth/security"
	"github.com/SomeHowMicroservice/shm-be/services/user/model"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authServiceImpl struct {
	cacheRepo  repository.CacheRepository
	userClient userpb.UserServiceClient
	mailer     smtp.Mailer
	cfg        *config.Config
	mqChannel  *amqp091.Channel
}

func NewAuthService(cacheRepo repository.CacheRepository, userClient userpb.UserServiceClient, mailer smtp.Mailer, cfg *config.Config, mqChannel *amqp091.Channel) AuthService {
	return &authServiceImpl{
		cacheRepo,
		userClient,
		mailer,
		cfg,
		mqChannel,
	}
}

func (s *authServiceImpl) SignUp(ctx context.Context, req *protobuf.SignUpRequest) (string, error) {
	// Kiểm tra Username và Email đã tồn tại chưa
	uRes, err := s.userClient.CheckUsernameExists(ctx, &userpb.CheckUsernameExistsRequest{Username: req.Username})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return "", fmt.Errorf("lỗi từ user service: %s", st.Message())
		}
		return "", fmt.Errorf("lỗi không xác định: %w", err)
	}
	if uRes.Exists {
		return "", customErr.ErrUsernameAlreadyExists
	}

	eRes, err := s.userClient.CheckEmailExists(ctx, &userpb.CheckEmailExistsRequest{Email: req.Email})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return "", fmt.Errorf("lỗi từ user service: %s", st.Message())
		}
		return "", fmt.Errorf("lỗi không xác định: %w", err)
	}
	if eRes.Exists {
		return "", customErr.ErrEmailAlreadyExists
	}

	// Tạo OTP và Registration Token trả ra cho FE
	otp := common.GenerateOTP(6)
	registrationToken := uuid.NewString()

	// Băm mật khẩu
	hashedPassword, err := common.HashPassword(req.Password)
	if err != nil {
		return "", fmt.Errorf("băm mật khẩu thất bại: %w", err)
	}

	// Tạo object dữ liệu người dùng đăng ký lưu vào Redis
	registrationData := &common.RegistrationData{
		Email:    req.Email,
		Username: req.Username,
		Password: hashedPassword,
		Otp:      otp,
		Attempts: 0,
	}

	// Lưu vô Redis
	if err = s.cacheRepo.SaveRegistrationData(ctx, registrationToken, registrationData, 3*time.Minute); err != nil {
		return "", err
	}

	// Tạo struct dữ liệu cần đẩy vào RabbitMQ
	emailMsg := common.AuthEmailMessage{
		To:      req.Email,
		Subject: "Xác thực đăng ký tài khoản tại SomeHow",
		Otp:     otp,
	}

	// Chuyển dổi
	body, err := json.Marshal(emailMsg)
	if err != nil {
		return "", fmt.Errorf("chuyển đổi EmailMessage thất bại: %w", err)
	}

	// Publish lên RabbitMQ
	if err := mq.PublishMessage(s.mqChannel, "", "email.send", body); err != nil {
		return "", fmt.Errorf("publish email msg thất bại: %w", err)
	}

	return registrationToken, nil
}

func (s *authServiceImpl) VerifySignUp(ctx context.Context, req *protobuf.VerifySignUpRequest) (*protobuf.AuthResponse, string, time.Duration, string, time.Duration, error) {
	// Lấy dữ liệu người dùng đăng ký từ redis ra
	regData, err := s.cacheRepo.GetRegistrationData(ctx, req.RegistrationToken)
	if err != nil {
		return nil, "", 0, "", 0, err
	}
	if regData == nil {
		return nil, "", 0, "", 0, customErr.ErrAuthDataNotFound
	}

	// Thử quá 3 lần thì xóa dữ liệu khỏi redis và báo lỗi
	if regData.Attempts >= 3 {
		if err = s.cacheRepo.DeleteAuthData(ctx, "sign-up", req.RegistrationToken); err != nil {
			return nil, "", 0, "", 0, err
		}
		return nil, "", 0, "", 0, customErr.ErrTooManyAttempts
	}

	// Thêm số lần thử và lưu lại vào redis
	regData.Attempts++
	if err = s.cacheRepo.SaveRegistrationData(ctx, req.RegistrationToken, regData, 3*time.Minute); err != nil {
		return nil, "", 0, "", 0, err
	}

	// Kiểm tra OTP có khơp không
	if regData.Otp != req.Otp {
		return nil, "", 0, "", 0, customErr.ErrInvalidOTP
	}

	// Kiểm tra lại tồn tại Username và Email cho nó chắc lỡ trong quá xác thực OTP thì có tạo tài khoản ở chỗ khác rồi
	uRes, err := s.userClient.CheckUsernameExists(ctx, &userpb.CheckUsernameExistsRequest{Username: regData.Username})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return nil, "", 0, "", 0, fmt.Errorf("lỗi từ user service: %s", st.Message())
		}
		return nil, "", 0, "", 0, fmt.Errorf("lỗi không xác định: %w", err)
	}
	if uRes.Exists {
		return nil, "", 0, "", 0, customErr.ErrUsernameAlreadyExists
	}

	eRes, err := s.userClient.CheckEmailExists(ctx, &userpb.CheckEmailExistsRequest{Email: regData.Email})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return nil, "", 0, "", 0, fmt.Errorf("lỗi từ user service: %s", st.Message())
		}
		return nil, "", 0, "", 0, fmt.Errorf("lỗi không xác định: %w", err)
	}
	if eRes.Exists {
		return nil, "", 0, "", 0, customErr.ErrEmailAlreadyExists
	}

	// Tạo struct và gửi tới User Service bằng gRPC
	newUser := &userpb.CreateUserRequest{
		Username: regData.Username,
		Email:    regData.Email,
		Password: regData.Password,
	}
	userRes, err := s.userClient.CreateUser(ctx, newUser)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return nil, "", 0, "", 0, fmt.Errorf("lỗi từ user service: %s", st.Message())
		}
		return nil, "", 0, "", 0, fmt.Errorf("lỗi không xác định: %w", err)
	}

	// Tạo Access Token và Refresh Token
	accessToken, err := security.GenerateToken(userRes.Id, userRes.Roles, s.cfg.Jwt.SecretKey, s.cfg.Jwt.AccessExpiresIn)
	if err != nil {
		return nil, "", 0, "", 0, fmt.Errorf("tạo token xác thực thất bại")
	}
	refreshToken, err := security.GenerateToken(userRes.Id, userRes.Roles, s.cfg.Jwt.SecretKey, s.cfg.Jwt.RefreshExpiresIn)
	if err != nil {
		return nil, "", 0, "", 0, fmt.Errorf("tạo token xác thực thất bại")
	}

	// Xóa dữ liệu trong redis
	if err = s.cacheRepo.DeleteAuthData(ctx, "sign-up", req.RegistrationToken); err != nil {
		return nil, "", 0, "", 0, err
	}
	// Tạo struct phản hồi
	authRes := &protobuf.AuthResponse{
		Id:        userRes.Id,
		Username:  userRes.Username,
		Email:     userRes.Email,
		CreatedAt: userRes.CreatedAt,
		Profile: &protobuf.ProfileResponse{
			Id:        userRes.Profile.Id,
			FirstName: userRes.Profile.FirstName,
			LastName:  userRes.Profile.LastName,
			Gender:    userRes.Profile.Gender,
			Dob:       userRes.Profile.Dob,
		},
	}

	return authRes, accessToken, s.cfg.Jwt.AccessExpiresIn, refreshToken, s.cfg.Jwt.RefreshExpiresIn, nil
}

func (s *authServiceImpl) SignIn(ctx context.Context, req *protobuf.SignInRequest) (*protobuf.AuthResponse, string, time.Duration, string, time.Duration, error) {
	// Tìm người dùng theo username
	userRes, err := s.userClient.GetUserByUsername(ctx, &userpb.GetUserByUsernameRequest{Username: req.Username})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, "", 0, "", 0, customErr.ErrUserNotFound
			default:
				return nil, "", 0, "", 0, fmt.Errorf("lỗi từ user service: %s", st.Message())
			}
		}
		return nil, "", 0, "", 0, fmt.Errorf("lỗi không xác định: %w", err)
	}
	// Kiểm tra mật khẩu
	isCorrectPassword := common.VerifyPassword(req.Password, userRes.Password)
	if !isCorrectPassword {
		return nil, "", 0, "", 0, customErr.ErrInvalidPassword
	}

	// Kiểm tra quyền
	if !hasRoleUser(model.RoleUser, userRes.Roles) {
		return nil, "", 0, "", 0, customErr.ErrForbidden
	}

	// Tạo Access Token và Refresh Token
	accessToken, err := security.GenerateToken(userRes.Id, userRes.Roles, s.cfg.Jwt.SecretKey, s.cfg.Jwt.AccessExpiresIn)
	if err != nil {
		return nil, "", 0, "", 0, fmt.Errorf("tạo token xác thực thất bại")
	}
	refreshToken, err := security.GenerateToken(userRes.Id, userRes.Roles, s.cfg.Jwt.SecretKey, s.cfg.Jwt.RefreshExpiresIn)
	if err != nil {
		return nil, "", 0, "", 0, fmt.Errorf("tạo token xác thực thất bại")
	}
	authRes := &protobuf.AuthResponse{
		Id:        userRes.Id,
		Username:  userRes.Username,
		Email:     userRes.Email,
		CreatedAt: userRes.CreatedAt,
		Profile: &protobuf.ProfileResponse{
			Id:        userRes.Profile.Id,
			FirstName: userRes.Profile.FirstName,
			LastName:  userRes.Profile.LastName,
			Gender:    userRes.Profile.Gender,
			Dob:       userRes.Profile.Dob,
		},
	}
	return authRes, accessToken, s.cfg.Jwt.AccessExpiresIn, refreshToken, s.cfg.Jwt.RefreshExpiresIn, nil
}

func (s *authServiceImpl) RefreshToken(ctx context.Context, req *protobuf.RefreshTokenRequest) (string, time.Duration, string, time.Duration, error) {
	accessToken, err := security.GenerateToken(req.Id, req.Roles, s.cfg.Jwt.SecretKey, s.cfg.Jwt.AccessExpiresIn)
	if err != nil {
		return "", 0, "", 0, fmt.Errorf("tạo token xác thực thất bại")
	}

	refreshToken, err := security.GenerateToken(req.Id, req.Roles, s.cfg.Jwt.SecretKey, s.cfg.Jwt.RefreshExpiresIn)
	if err != nil {
		return "", 0, "", 0, fmt.Errorf("tạo token xác thực thất bại")
	}

	return accessToken, s.cfg.Jwt.AccessExpiresIn, refreshToken, s.cfg.Jwt.RefreshExpiresIn, nil
}

func (s *authServiceImpl) ChangePassword(ctx context.Context, req *protobuf.ChangePasswordRequest) (string, time.Duration, string, time.Duration, error) {
	// Lấy user theo id có mật khẩu
	userRes, err := s.userClient.GetUserById(ctx, &userpb.GetUserByIdRequest{Id: req.Id})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return "", 0, "", 0, customErr.ErrUserNotFound
			default:
				return "", 0, "", 0, fmt.Errorf("lỗi từ user service: %s", st.Message())
			}
		}
		return "", 0, "", 0, fmt.Errorf("lỗi không xác định: %w", err)
	}
	// Kiểm tra mật khẩu
	isCorrectPassword := common.VerifyPassword(req.OldPassword, userRes.Password)
	if !isCorrectPassword {
		return "", 0, "", 0, customErr.ErrInvalidPassword
	}

	hashedNewPassword, err := common.HashPassword(req.NewPassword)
	if err != nil {
		return "", 0, "", 0, fmt.Errorf("băm mật khẩu thất bại: %w", err)
	}

	_, err = s.userClient.UpdateUserPassword(ctx, &userpb.UpdateUserPasswordRequest{
		Id:          userRes.Id,
		NewPassword: hashedNewPassword,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return "", 0, "", 0, customErr.ErrUserNotFound
			default:
				return "", 0, "", 0, fmt.Errorf("lỗi từ user service: %s", st.Message())
			}
		}
		return "", 0, "", 0, fmt.Errorf("lỗi không xác định: %w", err)
	}
	accessToken, err := security.GenerateToken(userRes.Id, userRes.Roles, s.cfg.Jwt.SecretKey, s.cfg.Jwt.AccessExpiresIn)
	if err != nil {
		return "", 0, "", 0, fmt.Errorf("tạo token xác thực thất bại")
	}
	refreshToken, err := security.GenerateToken(userRes.Id, userRes.Roles, s.cfg.Jwt.SecretKey, s.cfg.Jwt.RefreshExpiresIn)
	if err != nil {
		return "", 0, "", 0, fmt.Errorf("tạo token xác thực thất bại")
	}
	return accessToken, s.cfg.Jwt.AccessExpiresIn, refreshToken, s.cfg.Jwt.RefreshExpiresIn, nil
}

func (s *authServiceImpl) ForgotPassword(ctx context.Context, req *protobuf.ForgotPasswordRequest) (string, error) {
	// Kiểm tra email có tồn tại không
	res, err := s.userClient.CheckEmailExists(ctx, &userpb.CheckEmailExistsRequest{
		Email: req.Email,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return "", fmt.Errorf("lỗi từ user service: %s", st.Message())
		}
		return "", fmt.Errorf("lỗi không xác định: %w", err)
	}
	if !res.Exists {
		return "", customErr.ErrUserNotFound
	}

	// Tạo OTP và Forgot Password Token trả ra cho FE
	otp := common.GenerateOTP(6)
	forgotPasswordToken := uuid.NewString()

	// Tạo struct dữ liệu quên mật khẩu
	forgotPasswordData := &common.ForgotPasswordData{
		Email:    req.Email,
		Otp:      otp,
		Attempts: 0,
	}

	// Lưu vào redis
	if err = s.cacheRepo.SaveForgotPasswordData(ctx, forgotPasswordToken, forgotPasswordData, 3*time.Minute); err != nil {
		return "", err
	}

	// Tạo struct dữ liệu cần đẩy vào RabbitMQ
	emailMsg := common.AuthEmailMessage{
		To:      req.Email,
		Subject: "Xác thực quên mật khẩu tại SomeHow",
		Otp:     otp,
	}

	// Chuyển dổi sang json
	body, err := json.Marshal(emailMsg)
	if err != nil {
		return "", fmt.Errorf("chuyển đổi EmailMessage thất bại: %w", err)
	}

	// Publish lên RabbitMQ
	if err := mq.PublishMessage(s.mqChannel, "", "email.send", body); err != nil {
		return "", fmt.Errorf("publish email msg thất bại: %w", err)
	}

	return forgotPasswordToken, nil
}

func (s *authServiceImpl) VerifyForgotPassword(ctx context.Context, req *protobuf.VerifyForgotPasswordRequest) (string, error) {
	// Lấy dữ liệu quên mạt khẩu từ Redis
	forgData, err := s.cacheRepo.GetForgotPasswordData(ctx, req.ForgotPasswordToken)
	if err != nil {
		return "", err
	}
	if forgData == nil {
		return "", customErr.ErrAuthDataNotFound
	}

	// Thử quá 3 lần thì xóa dữ liệu khỏi redis và báo lỗi
	if forgData.Attempts >= 3 {
		if err = s.cacheRepo.DeleteAuthData(ctx, "forgot-password", req.ForgotPasswordToken); err != nil {
			return "", err
		}
		return "", customErr.ErrTooManyAttempts
	}

	if forgData.Otp != req.Otp {
		return "", customErr.ErrInvalidOTP
	}

	resetPasswordToken := uuid.NewString()

	if err = s.cacheRepo.SaveResetPasswordData(ctx, resetPasswordToken, forgData.Email, 3*time.Minute); err != nil {
		return "", err
	}

	if err = s.cacheRepo.DeleteAuthData(ctx, "forgot-password", req.ForgotPasswordToken); err != nil {
		return "", err
	}

	return resetPasswordToken, nil
}

func (s *authServiceImpl) ResetPassword(ctx context.Context, req *protobuf.ResetPasswordRequest) error {
	// Lấy dữ liệu làm mới mật khẩu từ Redis
	email, err := s.cacheRepo.GetResetPasswordData(ctx, req.ResetPasswordToken)
	if err != nil {
		return err
	}
	if email == "" {
		return customErr.ErrAuthDataNotFound
	}

	userRes, err := s.userClient.GetUserPublicByEmail(ctx, &userpb.GetUserByEmailRequest{
		Email: email,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return customErr.ErrUserNotFound
			default:
				return fmt.Errorf("lỗi từ user service: %s", st.Message())
			}
		}
		return fmt.Errorf("lỗi không xác định: %w", err)
	}

	hashedNewPassword, err := common.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("băm mật khẩu thất bại: %w", err)
	}

	_, err = s.userClient.UpdateUserPassword(ctx, &userpb.UpdateUserPasswordRequest{
		Id:          userRes.Id,
		NewPassword: hashedNewPassword,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return customErr.ErrUserNotFound
			default:
				return fmt.Errorf("lỗi từ user service: %s", st.Message())
			}
		}
		return fmt.Errorf("lỗi không xác định: %w", err)
	}

	return nil
}

func (s *authServiceImpl) AdminSignIn(ctx context.Context, req *protobuf.SignInRequest) (*protobuf.AuthResponse, string, time.Duration, string, time.Duration, error) {
	// Tìm người dùng theo username
	userRes, err := s.userClient.GetUserByUsername(ctx, &userpb.GetUserByUsernameRequest{Username: req.Username})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, "", 0, "", 0, customErr.ErrUserNotFound
			default:
				return nil, "", 0, "", 0, fmt.Errorf("lỗi từ user service: %s", st.Message())
			}
		}
		return nil, "", 0, "", 0, fmt.Errorf("lỗi không xác định: %w", err)
	}
	// Kiểm tra mật khẩu
	isCorrectPassword := common.VerifyPassword(req.Password, userRes.Password)
	if !isCorrectPassword {
		return nil, "", 0, "", 0, customErr.ErrInvalidPassword
	}

	// Kiểm tra quyền
	if !isAdmin(model.RoleUser, userRes.Roles) {
		return nil, "", 0, "", 0, customErr.ErrForbidden
	}

	// Tạo Access Token và Refresh Token
	accessToken, err := security.GenerateToken(userRes.Id, userRes.Roles, s.cfg.Jwt.SecretKey, s.cfg.Jwt.AccessExpiresIn)
	if err != nil {
		return nil, "", 0, "", 0, fmt.Errorf("tạo token xác thực thất bại")
	}
	refreshToken, err := security.GenerateToken(userRes.Id, userRes.Roles, s.cfg.Jwt.SecretKey, s.cfg.Jwt.RefreshExpiresIn)
	if err != nil {
		return nil, "", 0, "", 0, fmt.Errorf("tạo token xác thực thất bại")
	}
	authRes := &protobuf.AuthResponse{
		Id:        userRes.Id,
		Username:  userRes.Username,
		Email:     userRes.Email,
		CreatedAt: userRes.CreatedAt,
		Profile: &protobuf.ProfileResponse{
			Id:        userRes.Profile.Id,
			FirstName: userRes.Profile.FirstName,
			LastName:  userRes.Profile.LastName,
			Gender:    userRes.Profile.Gender,
			Dob:       userRes.Profile.Dob,
		},
	}
	return authRes, accessToken, s.cfg.Jwt.AccessExpiresIn, refreshToken, s.cfg.Jwt.RefreshExpiresIn, nil
}

func hasRoleUser(role string, roles []string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}

	return false
}

func isAdmin(role string, roles []string) bool {
	if len(roles) > 1 && hasRoleUser(role, roles) {
		return true
	}

	return false
}
