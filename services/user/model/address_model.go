package model

import "time"

type Address struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
	FullName    string    `gorm:"type:varchar(100);not null" json:"full_name"`
	PhoneNumber string    `gorm:"type:char(10);not null" json:"phone_number"`
	Street      string    `gorm:"type:varchar(100);not null" json:"street"`
	Ward        string    `gorm:"type:varchar(50);not null" json:"ward"`
	Province    string    `gorm:"type:varchar(50);not null" json:"province"`
	IsDefault   bool      `gorm:"type:boolean;not null" json:"is_default"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	UserID      string    `gorm:"type:char(36);not null" json:"-"`

	User *User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}
