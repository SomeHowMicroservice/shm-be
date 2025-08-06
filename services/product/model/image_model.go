package model

import "time"

type Image struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
	ProductID   string    `gorm:"type:char(36);not null" json:"-"`
	ColorID     string    `gorm:"type:char(36);not null" json:"-"`
	Url         string    `gorm:"type:varchar(255);not null" json:"url"`
	SortOrder   int       `gorm:"type:int;not null" json:"sort_order"`
	IsThumbnail bool      `gorm:"type:boolean;not null" json:"is_thumbnail"`
	IsDeleted   bool      `gorm:"type:boolean;default:false" json:"is_deleted"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedByID string    `gorm:"type:char(36);not null" json:"created_by_id"`
	UpdatedByID string    `gorm:"type:char(36);not null" json:"updated_by_id"`

	Product *Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"product"`
	Color   *Color   `gorm:"foreignKey:ColorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"color"`
}
