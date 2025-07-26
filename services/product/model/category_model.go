package model

import "time"

type Category struct {
	ID          string      `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string      `gorm:"type:varchar(100);not null" json:"name"`
	Slug        string      `gorm:"type:varchar(100);uniqueIndex:categories_slug_key;not null" json:"slug"`
	Parents     []*Category `gorm:"many2many:category_parents;joinForeignKey:ChildID;joinReferences:ParentID" json:"parents"`
	Children    []*Category `gorm:"many2many:category_parents;joinForeignKey:ParentID;joinReferences:ChildID" json:"children"`
	CreatedAt   time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedByID string      `gorm:"type:char(36);not null" json:"created_by_id"`
	UpdatedByID string      `gorm:"type:char(36);not null" json:"updated_by_id"`
}
