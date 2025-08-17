package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/post/model"
)

type TopicRepository interface {
	Create(ctx context.Context, topic *model.Topic) error
}