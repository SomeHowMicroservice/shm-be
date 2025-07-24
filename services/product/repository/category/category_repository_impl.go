package repository

import (
	"context"
	"errors"
	"time"

	// customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/product/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type categoryRepository struct {
	collection *mongo.Collection
}

func NewCategoryRepository(db *mongo.Database) CategoryRepository {
	collection := db.Collection("categories")
	return &categoryRepository{collection}
}

func (r *categoryRepository) Create(ctx context.Context, category *model.Category) error {
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	if category.ParentIDs == nil {
		category.ParentIDs = []bson.ObjectID{}
	}
	if category.ChildrenIDs == nil {
		category.ChildrenIDs = []bson.ObjectID{}
	}
	
	_, err := r.collection.InsertOne(ctx, category);
	if err != nil {
		return err
	}
	return nil
}

func (r *categoryRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"slug": slug})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *categoryRepository) ExistsByID(ctx context.Context, id bson.ObjectID) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *categoryRepository) FindByID(ctx context.Context, id bson.ObjectID) (*model.Category, error) {
	var category model.Category
	if err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&category); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}