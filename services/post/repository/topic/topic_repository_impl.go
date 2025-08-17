package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/post/model"
	"gorm.io/gorm"
)

type topicRepositoryImpl struct {
	db *gorm.DB
}

func NewTopicRepository(db *gorm.DB) TopicRepository {
	return &topicRepositoryImpl{db}
}

func (r *topicRepositoryImpl) Create(ctx context.Context, topic *model.Topic) error {
	if err := r.db.WithContext(ctx).Create(topic).Error; err != nil {
		return err
	}

	return nil
}