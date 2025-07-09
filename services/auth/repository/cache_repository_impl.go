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
	return &cacheRepositoryImpl{rdb}
}

const (
	serviceName = "auth-service"
)

func (r *cacheRepositoryImpl) SaveRegistrationData(ctx context.Context, registrationToken string, data common.RegistrationData, ttl time.Duration) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("chuyển đổi dữ liệu thất bại: %w", err)
	}
	redisKey := fmt.Sprintf("%s:sign-up:%s", serviceName, registrationToken)
	if err := r.rdb.Set(ctx, redisKey, bytes, ttl).Err(); err != nil {
		return fmt.Errorf("lưu dữ liệu thất bại: %w", err)
	}
	return nil
}

func (r *cacheRepositoryImpl) GetRegistrationData(ctx context.Context, registrationToken string) (*common.RegistrationData, error) {
	redisKey := fmt.Sprintf("%s:sign-up:%s", serviceName, registrationToken)
	regDataJSON, err := r.rdb.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("lấy dữ liệu từ bộ nhớ tạm thất bại: %w", err)
	}
	var regData common.RegistrationData
	if err = json.Unmarshal([]byte(regDataJSON), &regData); err != nil {
		return nil, fmt.Errorf("chuyển đổi dữ liệu thất bại: %w", err)
	}
	return &regData, nil
}

func (r *cacheRepositoryImpl) DeleteAuthData(ctx context.Context, keyType string, token string) error {
	redisKey := fmt.Sprintf("%s:%s:%s", serviceName, keyType, token) 
	if err := r.rdb.Del(ctx, redisKey).Err(); err != nil {
		return fmt.Errorf("xóa dữ liệu thất bại: %w", err)
	}
	return nil
}
