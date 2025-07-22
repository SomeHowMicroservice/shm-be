package service

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/user/model"
	"github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
)

type UserService interface {
	CheckEmailExists(ctx context.Context, email string) (bool, error)

	CheckUsernameExists(ctx context.Context, username string) (bool, error)

	CreateUser(ctx context.Context, req *protobuf.CreateUserRequest) (*model.User, error)

	UpdateUserPassword(ctx context.Context, req *protobuf.UpdateUserPasswordRequest) error

	GetUserByUsername(ctx context.Context, username string) (*model.User, error)

	GetUserByID(ctx context.Context, id string) (*model.User, error)

	GetUserByEmail(ctx context.Context, email string) (*model.User, error)

	UpdateUserProfile(ctx context.Context, req *protobuf.UpdateUserProfileRequest) (*model.User, error)

	UpdateUserMeasurement(ctx context.Context, req *protobuf.UpdateUserMeasurementRequest) (*model.Measurement, error)

	GetMeasurementByUserID(ctx context.Context, userID string) (*model.Measurement, error)

	GetAddressesByUserID(ctx context.Context, userID string) ([]*model.Address, error)

	CreateAddress(ctx context.Context, req *protobuf.CreateAddressRequest) (*model.Address, error)

	UpdateAddress(ctx context.Context, req *protobuf.UpdateAddressRequest) (*model.Address, error)

	DeleteAddress(ctx context.Context, req *protobuf.DeleteAddressRequest) error
}
