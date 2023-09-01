package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-server/internal/domain"
	"github.com/huangyul/go-server/internal/repository"
	"github.com/huangyul/go-server/internal/repository/dao"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidUserOrPassword = errors.New("incorrect email or password")
var ErrNotFound = dao.ErrNotFound

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(userRep *repository.UserRepository) *UserService {
	return &UserService{
		repo: userRep,
	}
}

func (srv *UserService) SignUp(ctx *gin.Context, user *domain.User) error {
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("system error")
	}
	user.Password = string(password)
	return srv.repo.Create(ctx, user)
}

func (srv *UserService) FindByEmail(ctx *gin.Context, user domain.User) (string, error) {
	u, err := srv.repo.FindByEmail(ctx, user)
	if errors.Is(err, ErrNotFound) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		return "", ErrInvalidUserOrPassword
	}
	return "", nil
}
