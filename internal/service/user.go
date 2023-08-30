package service

import "github.com/huangyul/go-server/internal/repository"

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(userRep *repository.UserRepository) *UserService {
	return &UserService{
		repo: userRep,
	}
}
