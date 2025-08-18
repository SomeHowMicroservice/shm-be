package container

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	postpb "github.com/SomeHowMicroservice/shm-be/services/post/protobuf"
)

type PostContainer struct {
	Handler *handler.PostHandler
}

func NewPostContainer(postClient postpb.PostServiceClient) *PostContainer {
	handler := handler.NewPostHandler(postClient)
	return &PostContainer{handler}
}
