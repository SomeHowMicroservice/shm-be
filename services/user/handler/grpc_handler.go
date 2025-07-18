package handler

import (
	"context"
	"time"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/user/model"
	"github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"github.com/SomeHowMicroservice/shm-be/services/user/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCHandler struct {
	protobuf.UnimplementedUserServiceServer
	svc service.UserService
}

func NewGRPCHandler(grpcServer *grpc.Server, svc service.UserService) *GRPCHandler {
	return &GRPCHandler{svc: svc}
}

func (h *GRPCHandler) CheckEmailExists(ctx context.Context, req *protobuf.CheckEmailExistsRequest) (*protobuf.UserCheckedResponse, error) {
	exists, err := h.svc.CheckEmailExists(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &protobuf.UserCheckedResponse{
		Exists: exists,
	}, nil
}

func (h *GRPCHandler) CheckUsernameExists(ctx context.Context, req *protobuf.CheckUsernameExistsRequest) (*protobuf.UserCheckedResponse, error) {
	exists, err := h.svc.CheckUsernameExists(ctx, req.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &protobuf.UserCheckedResponse{
		Exists: exists,
	}, nil
}

func (h *GRPCHandler) CreateUser(ctx context.Context, req *protobuf.CreateUserRequest) (*protobuf.UserPublicResponse, error) {
	user, err := h.svc.CreateUser(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrUsernameAlreadyExists, customErr.ErrEmailAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return toUserPublicResponse(user), nil
}

func (h *GRPCHandler) GetUserByUsername(ctx context.Context, req *protobuf.GetUserByUsernameRequest) (*protobuf.UserResponse, error) {
	user, err := h.svc.GetUserByUsername(ctx, req.Username)
	if err != nil {
		switch err {
		case customErr.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return toUserResponse(user), nil
}

func (h *GRPCHandler) GetUserPublicById(ctx context.Context, req *protobuf.GetUserByIdRequest) (*protobuf.UserPublicResponse, error) {
	user, err := h.svc.GetUserByID(ctx, req.Id)
	if err != nil {
		switch err {
		case customErr.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return toUserPublicResponse(user), nil
}

func (h *GRPCHandler) GetUserPublicByEmail(ctx context.Context, req *protobuf.GetUserByEmailRequest) (*protobuf.UserPublicResponse, error) {
	user, err := h.svc.GetUserByEmail(ctx, req.Email)
	if err != nil {
		switch err {
		case customErr.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return toUserPublicResponse(user), nil
}

func (h *GRPCHandler) GetUserById(ctx context.Context, req *protobuf.GetUserByIdRequest) (*protobuf.UserResponse, error) {
	user, err := h.svc.GetUserByID(ctx, req.Id)
	if err != nil {
		switch err {
		case customErr.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return toUserResponse(user), nil
}

func (h *GRPCHandler) UpdateUserPassword(ctx context.Context, req *protobuf.UpdateUserPasswordRequest) (*protobuf.UserUpdatedResponse, error) {
	if err := h.svc.UpdateUserPassword(ctx, req); err != nil {
		switch err {
		case customErr.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.UserUpdatedResponse{
		Success: true,
	}, nil
}

func (h *GRPCHandler) UpdateUserProfile(ctx context.Context, req *protobuf.UpdateUserProfileRequest) (*protobuf.UserPublicResponse, error) {
	user, err := h.svc.UpdateUserProfile(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrProfileNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return toUserPublicResponse(user), nil
}

func (h *GRPCHandler) GetMeasurementByUserId(ctx context.Context, req *protobuf.GetMeasurementByUserIdRequest) (*protobuf.MeasurementResponse, error) {
	measurement, err := h.svc.GetMeasurementByUserID(ctx, req.UserId)
	if err != nil {
		switch err {
		case customErr.ErrMeasurementNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return toMeasurementResponse(measurement), nil
}

func (h *GRPCHandler) UpdateUserMeasurement(ctx context.Context, req *protobuf.UpdateUserMeasurementRequest) (*protobuf.MeasurementResponse, error) {
	measurement, err := h.svc.UpdateUserMeasurement(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrMeasurementNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case customErr.ErrForbidden:
			return nil, status.Error(codes.PermissionDenied, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return toMeasurementResponse(measurement), nil
}

func toMeasurementResponse(measurement *model.Measurement) *protobuf.MeasurementResponse {
	return &protobuf.MeasurementResponse{
		Id:     measurement.ID,
		Height: int32(measurement.Height),
		Weight: int32(measurement.Weight),
		Chest:  int32(measurement.Chest),
		Waist:  int32(measurement.Waist),
		Butt:   int32(measurement.Butt),
	}
}

func toUserResponse(user *model.User) *protobuf.UserResponse {
	roles := []string{}
	for _, r := range user.Roles {
		roles = append(roles, r.Name)
	}
	var dob string
	if user.Profile.DOB != nil {
		dob = user.Profile.DOB.Format("2006-01-02")
	}
	var gender string
	if user.Profile.Gender != nil {
		gender = *user.Profile.Gender
	}
	return &protobuf.UserResponse{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Roles:     roles,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		Profile: &protobuf.ProfileResponse{
			Id:        user.Profile.ID,
			FirstName: user.Profile.FirstName,
			LastName:  user.Profile.LastName,
			Gender:    gender,
			Dob:       dob,
		},
	}
}

func toUserPublicResponse(user *model.User) *protobuf.UserPublicResponse {
	roles := []string{}
	for _, r := range user.Roles {
		roles = append(roles, r.Name)
	}
	var dob string
	if user.Profile.DOB != nil {
		dob = user.Profile.DOB.Format("2006-01-02")
	}
	var gender string
	if user.Profile.Gender != nil {
		gender = *user.Profile.Gender
	}
	return &protobuf.UserPublicResponse{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Roles:     roles,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		Profile: &protobuf.ProfileResponse{
			Id:        user.Profile.ID,
			FirstName: user.Profile.FirstName,
			LastName:  user.Profile.LastName,
			Gender:    gender,
			Dob:       dob,
		},
	}
}
