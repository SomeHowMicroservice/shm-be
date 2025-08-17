package service

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/post/protobuf"
)

type PostService interface {
	CreateTopic(ctx context.Context, req *protobuf.CreateTopicRequest) (string, error)
}