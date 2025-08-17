package service

import (
	"context"
	"errors"
	"fmt"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/post/common"
	"github.com/SomeHowMicroservice/shm-be/services/post/model"
	"github.com/SomeHowMicroservice/shm-be/services/post/protobuf"
	topicRepo "github.com/SomeHowMicroservice/shm-be/services/post/repository/topic"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type postServiceImpl struct {
	topicRepo topicRepo.TopicRepository
}

func NewPostService(topicRepo topicRepo.TopicRepository) PostService {
	return &postServiceImpl{
		topicRepo,
	}
}

func (s *postServiceImpl) CreateTopic(ctx context.Context, req *protobuf.CreateTopicRequest) (string, error) {
	if req.Slug == nil {
		slug := common.GenerateSlug(*req.Slug)
		req.Slug = &slug
	}

	topic := &model.Topic{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Slug:        *req.Slug,
		CreatedByID: req.UserId,
		UpdatedByID: req.UserId,
	}
	if err := s.topicRepo.Create(ctx, topic); err != nil {
		if isUniqueViolation(err) {
			return "", customErr.ErrTopicAlreadyExists
		}
		return "", fmt.Errorf("tạo chủ đề bài viết thất bại: %w", err)
	}

	return topic.ID, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
