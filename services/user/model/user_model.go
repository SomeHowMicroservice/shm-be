package model

import "time"

type User struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	Username  string    `gorm:"type:varchar(50);uniqueIndex:users_username_key;not null" json:"username"`
	Email     string    `gorm:"type: varchar(150);uniqueIndex:users_email_key;not null" json:"email"`
	Password  string    `gorm:"type: varchar(255); not null" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Profile     *Profile     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"profile,omitempty"`
	Measurement *Measurement `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"measurement,omitempty"`
	Roles       []*Role      `gorm:"many2many:user_roles" json:"roles"`
	Addresses   []*Address   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"address,omitempty"`
}
