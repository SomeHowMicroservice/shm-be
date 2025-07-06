package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	Username  string    `gorm:"type:varchar(50);uniqueIndex:users_username_key;not null" json:"username"`
	Email     string    `gorm:"type: varchar(150);uniqueIndex:users_email_key;not null" json:"email"`
	Password  string    `gorm:"type: varchar(255); not null" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}