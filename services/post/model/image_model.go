package model

type Image struct {
	ID        string `gorm:"type:char(36);primaryKey" json:"id"`
	PostID    string `gorm:"type:char(36);not null" json:"-"`
	Url       string `gorm:"type:varchar(255);not null" json:"url"`
	FileID    string `gorm:"type:char(24)" json:"file_id"`
	SortOrder int    `gorm:"type:int;not null" json:"sort_order"`

	Post *Post `gorm:"foreignKey:PostID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"post"`
}
