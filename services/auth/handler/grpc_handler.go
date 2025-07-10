package handler

import (
	"context"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
	"github.com/SomeHowMicroservice/shm-be/services/auth/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcHandler struct {
	protobuf.UnimplementedAuthServiceServer
	svc service.AuthService
}

func NewGRPCHandler(grpcServer *grpc.Server, svc service.AuthService) *grpcHandler {
	return &grpcHandler{svc: svc}
}

func (h *grpcHandler) SignUp(ctx context.Context, req *protobuf.SignUpRequest) (*protobuf.SignUpResponse, error) {
	token, err := h.svc.SignUp(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrUsernameAlreadyExists, customErr.ErrEmailAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &protobuf.SignUpResponse{
		RegistrationToken: token,
	}, nil
}

func (h *grpcHandler) VerifySignUp(ctx context.Context, req *protobuf.VerifySignUpRequest) (*protobuf.VerifySignUpResponse, error) {
	user, accessToken, accessExpiresIn, refreshToken, refreshExpiresIn, err := h.svc.VerifySignUp(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrUsernameAlreadyExists, customErr.ErrEmailAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case customErr.ErrAuthDataNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case customErr.ErrTooManyAttempts, customErr.ErrInvalidOTP:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &protobuf.VerifySignUpResponse{
		User:             user,
		AccessToken:      accessToken,
		AccessExpiresIn:  int64(accessExpiresIn.Seconds()),
		RefreshToken:     refreshToken,
		RefreshExpiresIn: int64(refreshExpiresIn.Seconds()),
	}, nil
}

func (h *grpcHandler) SignIn(ctx context.Context, req *protobuf.SignInRequest) (*protobuf.SignInResponse, error) {
	user, accessToken, accessExpiresIn, refreshToken, refreshExpiresIn, err := h.svc.SignIn(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case customErr.ErrInvalidPassword:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &protobuf.SignInResponse{
		User:             user,
		AccessToken:      accessToken,
		AccessExpiresIn:  int64(accessExpiresIn.Seconds()),
		RefreshToken:     refreshToken,
		RefreshExpiresIn: int64(refreshExpiresIn.Seconds()),
	}, nil
}
