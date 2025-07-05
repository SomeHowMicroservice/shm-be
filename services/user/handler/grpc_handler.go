package handler

import "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"

type UserHandler struct {
	protobuf.UnimplementedUserServiceServer
	
}