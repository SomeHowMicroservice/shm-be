package handler

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/auth/common"
	protobuf "github.com/SomeHowMicroservice/shm-be/auth/protobuf"
	"github.com/SomeHowMicroservice/shm-be/auth/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCHandler struct {
	protobuf.UnimplementedAuthServiceServer
	svc service.AuthService
}

func NewGRPCHandler(grpcServer *grpc.Server, svc service.AuthService) *GRPCHandler {
	return &GRPCHandler{svc: svc}
}

func (h *GRPCHandler) SignUp(ctx context.Context, req *protobuf.SignUpRequest) (*protobuf.SignUpResponse, error) {
	token, err := h.svc.SignUp(ctx, req)
	if err != nil {
		switch err {
		case common.ErrUsernameAlreadyExists, common.ErrEmailAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.SignUpResponse{
		RegistrationToken: token,
	}, nil
}

func (h *GRPCHandler) VerifySignUp(ctx context.Context, req *protobuf.VerifySignUpRequest) (*protobuf.LoggedInResponse, error) {
	user, accessToken, accessExpiresIn, refreshToken, refreshExpiresIn, err := h.svc.VerifySignUp(ctx, req)
	if err != nil {
		switch err {
		case common.ErrUsernameAlreadyExists, common.ErrEmailAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case common.ErrAuthDataNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case common.ErrTooManyAttempts, common.ErrInvalidOTP:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &protobuf.LoggedInResponse{
		User:             user,
		AccessToken:      accessToken,
		AccessExpiresIn:  int64(accessExpiresIn.Seconds()),
		RefreshToken:     refreshToken,
		RefreshExpiresIn: int64(refreshExpiresIn.Seconds()),
	}, nil
}

func (h *GRPCHandler) SignIn(ctx context.Context, req *protobuf.SignInRequest) (*protobuf.LoggedInResponse, error) {
	user, accessToken, accessExpiresIn, refreshToken, refreshExpiresIn, err := h.svc.SignIn(ctx, req)
	if err != nil {
		switch err {
		case common.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case common.ErrForbidden:
			return nil, status.Error(codes.PermissionDenied, err.Error())
		case common.ErrInvalidPassword:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.LoggedInResponse{
		User:             user,
		AccessToken:      accessToken,
		AccessExpiresIn:  int64(accessExpiresIn.Seconds()),
		RefreshToken:     refreshToken,
		RefreshExpiresIn: int64(refreshExpiresIn.Seconds()),
	}, nil
}

func (h *GRPCHandler) RefreshToken(ctx context.Context, req *protobuf.RefreshTokenRequest) (*protobuf.RefreshTokenResponse, error) {
	accessToken, accessExpiresIn, refreshToken, refreshExpiresIn, err := h.svc.RefreshToken(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &protobuf.RefreshTokenResponse{
		AccessToken:      accessToken,
		AccessExpiresIn:  int64(accessExpiresIn.Seconds()),
		RefreshToken:     refreshToken,
		RefreshExpiresIn: int64(refreshExpiresIn.Seconds()),
	}, nil
}

func (h *GRPCHandler) ChangePassword(ctx context.Context, req *protobuf.ChangePasswordRequest) (*protobuf.RefreshTokenResponse, error) {
	accessToken, accessExpiresIn, refreshToken, refreshExpiresIn, err := h.svc.ChangePassword(ctx, req)
	if err != nil {
		switch err {
		case common.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case common.ErrInvalidPassword:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &protobuf.RefreshTokenResponse{
		AccessToken:      accessToken,
		AccessExpiresIn:  int64(accessExpiresIn.Seconds()),
		RefreshToken:     refreshToken,
		RefreshExpiresIn: int64(refreshExpiresIn.Seconds()),
	}, nil
}

func (h *GRPCHandler) ForgotPassword(ctx context.Context, req *protobuf.ForgotPasswordRequest) (*protobuf.ForgotPasswordResponse, error) {
	token, err := h.svc.ForgotPassword(ctx, req)
	if err != nil {
		switch err {
		case common.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.ForgotPasswordResponse{
		ForgotPasswordToken: token,
	}, nil
}

func (h *GRPCHandler) VerifyForgotPassword(ctx context.Context, req *protobuf.VerifyForgotPasswordRequest) (*protobuf.VerifyForgotPasswordResponse, error) {
	token, err := h.svc.VerifyForgotPassword(ctx, req)
	if err != nil {
		switch err {
		case common.ErrAuthDataNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case common.ErrTooManyAttempts, common.ErrInvalidOTP:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.VerifyForgotPasswordResponse{
		ResetPasswordToken: token,
	}, nil
}

func (h *GRPCHandler) ResetPassword(ctx context.Context, req *protobuf.ResetPasswordRequest) (*protobuf.UpdatedResponse, error) {
	if err := h.svc.ResetPassword(ctx, req); err != nil {
		switch err {
		case common.ErrUserNotFound, common.ErrAuthDataNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.UpdatedResponse{
		Success: true,
	}, nil
}

func (h *GRPCHandler) AdminSignIn(ctx context.Context, req *protobuf.SignInRequest) (*protobuf.LoggedInResponse, error) {
	user, accessToken, accessExpiresIn, refreshToken, refreshExpiresIn, err := h.svc.AdminSignIn(ctx, req)
	if err != nil {
		switch err {
		case common.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case common.ErrForbidden:
			return nil, status.Error(codes.PermissionDenied, err.Error())
		case common.ErrInvalidPassword:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.LoggedInResponse{
		User:             user,
		AccessToken:      accessToken,
		AccessExpiresIn:  int64(accessExpiresIn.Seconds()),
		RefreshToken:     refreshToken,
		RefreshExpiresIn: int64(refreshExpiresIn.Seconds()),
	}, nil
}
