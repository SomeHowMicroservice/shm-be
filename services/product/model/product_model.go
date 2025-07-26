package model

import (
	"time"

	"github.com/SomeHowMicroservice/shm-be/services/user/model"
)

type Product struct {
	ID          string     `gorm:"type:char(36);primaryKey" json:"id"`
	Title       string     `gorm:"type:varchar(255);not null" json:"title"`
	Slug        string     `gorm:"type:varchar(255);uniqueIndex:products_slug_key;not null" json:"slug"`
	Description string     `gorm:"type:text" json:"description"`
	Price       float32    `gorm:"not null" json:"price"`
	IsSale      bool       `gorm:"not null" json:"is_sale"`
	SalePrice   *float32   `json:"sale_price"`
	StartSale   *time.Time `json:"start_sale"`
	EndSale     *time.Time `json:"end_sale"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedByID string     `gorm:"type:char(36);not null" json:"-"`
	UpdatedByID string     `gorm:"type:char(36);not null" json:"-"`

	Categories []*Category `gorm:"many2many:product_categories" json:"categories"`
	CreatedBy *model.User `gorm:"foreignKey:CreatedByID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"created_by"`
	UpdatedBy *model.User `gorm:"foreignKey:UpdatedByID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"updated_by"`
}
