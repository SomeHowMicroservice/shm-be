package handler

import (
	"context"
	"time"

	"github.com/SomeHowMicroservice/shm-be/user/common"
	"github.com/SomeHowMicroservice/shm-be/user/model"
	userpb "github.com/SomeHowMicroservice/shm-be/user/protobuf/user"
	"github.com/SomeHowMicroservice/shm-be/user/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCHandler struct {
	userpb.UnimplementedUserServiceServer
	svc service.UserService
}

func NewGRPCHandler(grpcServer *grpc.Server, svc service.UserService) *GRPCHandler {
	return &GRPCHandler{svc: svc}
}

func (h *GRPCHandler) CheckEmailExists(ctx context.Context, req *userpb.CheckEmailExistsRequest) (*userpb.UserCheckedResponse, error) {
	exists, err := h.svc.CheckEmailExists(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.UserCheckedResponse{
		Exists: exists,
	}, nil
}

func (h *GRPCHandler) CheckUsernameExists(ctx context.Context, req *userpb.CheckUsernameExistsRequest) (*userpb.UserCheckedResponse, error) {
	exists, err := h.svc.CheckUsernameExists(ctx, req.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.UserCheckedResponse{
		Exists: exists,
	}, nil
}

func (h *GRPCHandler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.UserPublicResponse, error) {
	user, err := h.svc.CreateUser(ctx, req)
	if err != nil {
		switch err {
		case common.ErrUsernameAlreadyExists, common.ErrEmailAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return toUserPublicResponse(user), nil
}

func (h *GRPCHandler) GetUserByUsername(ctx context.Context, req *userpb.GetUserByUsernameRequest) (*userpb.UserResponse, error) {
	user, err := h.svc.GetUserByUsername(ctx, req.Username)
	if err != nil {
		switch err {
		case common.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return toUserResponse(user), nil
}

func (h *GRPCHandler) GetUserPublicById(ctx context.Context, req *userpb.GetUserByIdRequest) (*userpb.UserPublicResponse, error) {
	user, err := h.svc.GetUserByID(ctx, req.Id)
	if err != nil {
		switch err {
		case common.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return toUserPublicResponse(user), nil
}

func (h *GRPCHandler) GetUserPublicByEmail(ctx context.Context, req *userpb.GetUserByEmailRequest) (*userpb.UserPublicResponse, error) {
	user, err := h.svc.GetUserByEmail(ctx, req.Email)
	if err != nil {
		switch err {
		case common.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return toUserPublicResponse(user), nil
}

func (h *GRPCHandler) GetUserById(ctx context.Context, req *userpb.GetUserByIdRequest) (*userpb.UserResponse, error) {
	user, err := h.svc.GetUserByID(ctx, req.Id)
	if err != nil {
		switch err {
		case common.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return toUserResponse(user), nil
}

func (h *GRPCHandler) UpdateUserPassword(ctx context.Context, req *userpb.UpdateUserPasswordRequest) (*userpb.UserUpdatedResponse, error) {
	if err := h.svc.UpdateUserPassword(ctx, req); err != nil {
		switch err {
		case common.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &userpb.UserUpdatedResponse{
		Success: true,
	}, nil
}

func (h *GRPCHandler) UpdateUserProfile(ctx context.Context, req *userpb.UpdateUserProfileRequest) (*userpb.UserPublicResponse, error) {
	user, err := h.svc.UpdateUserProfile(ctx, req)
	if err != nil {
		switch err {
		case common.ErrProfileNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return toUserPublicResponse(user), nil
}

func (h *GRPCHandler) GetAddressesByUserId(ctx context.Context, req *userpb.GetAddressesByUserIdRequest) (*userpb.AddressesResponse, error) {
	addresses, err := h.svc.GetAddressesByUserID(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toAddressesResponse(addresses), nil
}

func (h *GRPCHandler) GetMeasurementByUserId(ctx context.Context, req *userpb.GetMeasurementByUserIdRequest) (*userpb.MeasurementResponse, error) {
	measurement, err := h.svc.GetMeasurementByUserID(ctx, req.UserId)
	if err != nil {
		switch err {
		case common.ErrMeasurementNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return toMeasurementResponse(measurement), nil
}

func (h *GRPCHandler) UpdateUserMeasurement(ctx context.Context, req *userpb.UpdateUserMeasurementRequest) (*userpb.MeasurementResponse, error) {
	measurement, err := h.svc.UpdateUserMeasurement(ctx, req)
	if err != nil {
		switch err {
		case common.ErrMeasurementNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case common.ErrForbidden:
			return nil, status.Error(codes.PermissionDenied, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return toMeasurementResponse(measurement), nil
}

func (h *GRPCHandler) CreateAddress(ctx context.Context, req *userpb.CreateAddressRequest) (*userpb.AddressResponse, error) {
	address, err := h.svc.CreateAddress(ctx, req)
	if err != nil {
		switch err {
		case common.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return toAddressResponse(address), nil

}

func (h *GRPCHandler) UpdateAddress(ctx context.Context, req *userpb.UpdateAddressRequest) (*userpb.AddressResponse, error) {
	address, err := h.svc.UpdateAddress(ctx, req)
	if err != nil {
		switch err {
		case common.ErrAddressNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case common.ErrForbidden:
			return nil, status.Error(codes.PermissionDenied, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return toAddressResponse(address), nil
}

func (h *GRPCHandler) DeleteAddress(ctx context.Context, req *userpb.DeleteAddressRequest) (*userpb.AddressDeletedResponse, error) {
	if err := h.svc.DeleteAddress(ctx, req); err != nil {
		switch err {
		case common.ErrAddressNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case common.ErrForbidden:
			return nil, status.Error(codes.PermissionDenied, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &userpb.AddressDeletedResponse{
		Success: true,
	}, nil
}

func (h *GRPCHandler) GetUsersById(ctx context.Context, req *userpb.GetUsersByIdRequest) (*userpb.UsersPublicResponse, error) {
	users, err := h.svc.GetUsersByID(ctx, req.Ids)
	if err != nil {
		switch err {
		case common.ErrHasUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return toUsersPublicResponse(users), nil
}

func toUsersPublicResponse(users []*model.User) *userpb.UsersPublicResponse {
	var userResponses []*userpb.UserPublicResponse
	for _, user := range users {
		userResponses = append(userResponses, toUserPublicResponse(user))
	}

	return &userpb.UsersPublicResponse{
		Users: userResponses,
	}
}

func toAddressResponse(address *model.Address) *userpb.AddressResponse {
	return &userpb.AddressResponse{
		Id:          address.ID,
		FullName:    address.FullName,
		PhoneNumber: address.PhoneNumber,
		Street:      address.Street,
		Ward:        address.Ward,
		Province:    address.Province,
		IsDefault:   &address.IsDefault,
	}
}

func toAddressesResponse(addresses []*model.Address) *userpb.AddressesResponse {
	var addressResponses []*userpb.AddressResponse
	for _, addr := range addresses {
		addressResponses = append(addressResponses, toAddressResponse(addr))
	}

	return &userpb.AddressesResponse{
		Addresses: addressResponses,
	}
}

func toMeasurementResponse(measurement *model.Measurement) *userpb.MeasurementResponse {
	return &userpb.MeasurementResponse{
		Id:     measurement.ID,
		Height: int32(measurement.Height),
		Weight: int32(measurement.Weight),
		Chest:  int32(measurement.Chest),
		Waist:  int32(measurement.Waist),
		Butt:   int32(measurement.Butt),
	}
}

func toUserResponse(user *model.User) *userpb.UserResponse {
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

	return &userpb.UserResponse{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Roles:     roles,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		Profile: &userpb.ProfileResponse{
			Id:        user.Profile.ID,
			FirstName: user.Profile.FirstName,
			LastName:  user.Profile.LastName,
			Gender:    gender,
			Dob:       dob,
		},
	}
}

func toUserPublicResponse(user *model.User) *userpb.UserPublicResponse {
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

	return &userpb.UserPublicResponse{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Roles:     roles,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		Profile: &userpb.ProfileResponse{
			Id:        user.Profile.ID,
			FirstName: user.Profile.FirstName,
			LastName:  user.Profile.LastName,
			Gender:    gender,
			Dob:       dob,
		},
	}
}
