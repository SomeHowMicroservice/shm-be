package model

import "time"

type Profile struct {
	ID        string     `gorm:"type:char(36);primaryKey" json:"id"`
	FirstName string     `gorm:"type:varchar(50)" json:"first_name"`
	LastName  string     `gorm:"type:varchar(50)" json:"last_name"`
	Gender    *string    `gorm:"type:enum_genders" json:"gender"`
	DOB       *time.Time `gorm:"type:date" json:"dob"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"-"`
	UserID    string     `gorm:"type:char(36);uniqueIndex:profiles_user_id_key;not null" json:"-"`

	User *User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}
