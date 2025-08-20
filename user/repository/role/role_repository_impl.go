package repository

import (
	"context"
	"errors"

	"github.com/SomeHowMicroservice/shm-be/user/model"
	"gorm.io/gorm"
)

type roleRepositoryImpl struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepositoryImpl{db}
}

func (r *roleRepositoryImpl) FindByName(ctx context.Context, name string) (*model.Role, error) {
	var role model.Role
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepositoryImpl) CreateUserRoles(ctx context.Context, userID string, roleID string) error {
	if err := r.db.WithContext(ctx).Exec(`INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)`, userID, roleID).Error; err != nil {
		return err
	}
	return nil
}