package model

import "time"

type Post struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Slug        string    `gorm:"type:varchar(255);uniqueIndex:posts_slug_key;not null" json:"slug"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	Viewed      int       `gorm:"type:bigint;not null;default:0" json:"viewed"`
	IsPublished bool      `gorm:"type:boolean;not null;default:true" json:"is_published"`
	IsDeleted   bool      `gorm:"type:boolean;not null;default:false" json:"is_deleted"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedByID string    `gorm:"type:char(36);not null" json:"created_by_id"`
	UpdatedByID string    `gorm:"type:char(36);not null" json:"updated_by_id"`
	TopicID     string    `gorm:"type:char(36);not null" json:"-"`

	Topic  *Topic   `gorm:"foreignKey:TopicID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"topic"`
	Images []*Image `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE" json:"images"`
}
