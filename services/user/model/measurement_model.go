package model

import "time"

type Measurement struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	Height    int       `gorm:"type:int" json:"height"`
	Weight    int       `gorm:"type:int" json:"weight"`
	Chest     int       `gorm:"type:int" json:"chest"`
	Waist     int       `gorm:"type:int" json:"waist"`
	Butt      int       `gorm:"type:int" json:"butt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`
	UserID    string    `gorm:"type:char(36);uniqueIndex:measurements_user_id_key;not null" json:"-"`

	User *User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}
