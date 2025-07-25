package model

import (
	"time"

	"github.com/SomeHowMicroservice/shm-be/services/user/model"
)

type Category struct {
	ID          string      `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string      `gorm:"type:varchar(100);not null" json:"name"`
	Slug        string      `gorm:"type:varchar(100);uniqueIndex:categories_slug_key;not null" json:"slug"`
	Parents     []*Category `gorm:"many2many:category_parents;joinForeignKey:ChildID;joinReferences:ParentID" json:"parents"`
	Children    []*Category `gorm:"many2many:category_parents;joinForeignKey:ParentID;joinReferences:ChildID" json:"children"`
	CreatedAt   time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedByID string      `gorm:"type:char(36);not null" json:"-"`
	UpdatedByID string      `gorm:"type:char(36);not null" json:"-"`

	CreatedBy *model.User `gorm:"foreignKey:CreatedByID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"created_by"`
	UpdatedBy *model.User `gorm:"foreignKey:UpdatedByID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"updated_by"`
}
