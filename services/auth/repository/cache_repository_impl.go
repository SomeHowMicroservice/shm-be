package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/SomeHowMicroservice/shm-be/services/auth/common"
	"github.com/redis/go-redis/v9"
)

type cacheRepositoryImpl struct {
	rdb *redis.Client
}

func NewCacheRepository(rdb *redis.Client) CacheRepository {
	return &cacheRepositoryImpl{
		rdb: rdb,
	}
}

const (
	serviceName = "auth-service"
)

func(r *cacheRepositoryImpl) SaveRegistrationData(ctx context.Context, registrationToken string, data common.RegistrationData, ttl time.Duration) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	redisKey := fmt.Sprintf("%s:sign-up:%s", serviceName, registrationToken)
	if err := r.rdb.Set(ctx, redisKey, bytes, ttl).Err(); err != nil {
		return err
	}

	return nil
}
