package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/user/model"
)

type MeasurementRepository interface {
	Create(ctx context.Context, measurement *model.Measurement) error

	Update(ctx context.Context, id string, updateData map[string]interface{}) error

	FindByID(ctx context.Context, id string) (*model.Measurement, error)

	FindByUserID(ctx context.Context, userID string) (*model.Measurement, error)
}