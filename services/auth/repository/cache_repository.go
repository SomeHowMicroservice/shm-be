package repository

import (
	"context"
	"time"

	"github.com/SomeHowMicroservice/shm-be/services/auth/common"
)

type CacheRepository interface {
	SaveRegistrationData(ctx context.Context, registrationToken string, data common.RegistrationData, ttl time.Duration) error
}