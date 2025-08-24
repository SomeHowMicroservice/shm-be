package model

import "time"

const (
	RoleAdmin       = "admin"
	RoleUser        = "user"
	RoleStaff       = "staff"
	RoleContributor = "contributor"
)

type Role struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string    `gorm:"type:enum_roles;uniqueIndex:roles_name_key;not null" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Users []*User `gorm:"many2many:user_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"users"`
}
