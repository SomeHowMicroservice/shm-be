package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Category struct {
	ID          bson.ObjectID   `bson:"_id" json:"id"`
	Name        string          `bson:"name" json:"name"`
	Slug        string          `bson:"slug" json:"slug"`
	ParentIDs   []bson.ObjectID `bson:"parent_ids,omitempty" json:"parent_ids"`
	ChildrenIDs []bson.ObjectID `bson:"children_ids,omitempty" json:"children_ids"`
	CreatedAt   time.Time       `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time       `bson:"updated_at" json:"updated_at"`
	CreatedByID   string          `bson:"created_by" json:"created_by"`
	UpdatedByID   string          `bson:"updated_by" json:"updated_by"`
}
