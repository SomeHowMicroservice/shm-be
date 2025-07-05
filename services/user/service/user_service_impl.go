package service

import "github.com/SomeHowMicroservice/shm-be/services/user/repository"

type userServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{
		repo: repo,
	}
}