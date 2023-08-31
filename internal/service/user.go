package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-server/internal/domain"
	"github.com/huangyul/go-server/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(userRep *repository.UserRepository) *UserService {
	return &UserService{
		repo: userRep,
	}
}

func (u *UserService) SignUp(ctx *gin.Context, user *domain.User) error {
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("系统错误")
	}
	user.Password = string(password)
	return u.repo.Create(ctx, user)
}
