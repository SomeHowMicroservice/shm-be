package model

import "time"

type Topic struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(150);not null" json:"name"`
	Slug        string    `gorm:"type:varchar(150);uniqueIndex:topics_slug_key;not null" json:"slug"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedByID string    `gorm:"type:char(36);not null" json:"created_by_id"`
	UpdatedByID string    `gorm:"type:char(36);not null" json:"updated_by_id"`

	Posts []*Post `gorm:"foreignKey:TopicID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE" json:"posts"`
}
