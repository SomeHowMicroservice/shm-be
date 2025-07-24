package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Product struct {
	ID          bson.ObjectID   `bson:"_id" json:"id"`
	Title       string          `bson:"title" json:"title"`
	Description string          `bson:"description" json:"description"`
	CategoryIDs []bson.ObjectID `bson:"category_ids" json:"category_ids"`
	CreatedAt   time.Time       `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time       `bson:"updated_at" json:"updated_at"`
}
