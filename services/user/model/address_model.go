package model

import "time"

type Address struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
	Street      string    `gorm:"type:varchar(100)" json:"street"`
	Ward        string    `gorm:"type:varchar(50)" json:"ward"`
	FullName    string    `gorm:"type:string" json:"full_name"`
	PhoneNumber string    `gorm:"type:char(10)" json:"phone_number"`
	Province    string    `gorm:"type:varchar(50)" json:"province"`
	IsDefault   bool      `gorm:"type:boolean;default:false" json:"is_default"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"-"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"-"`
	UserID      string    `gorm:"type:char(36);not null" json:"-"`

	User *User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}
