package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-server/internal/domain"
	"github.com/huangyul/go-server/internal/repository/dao"
)

var ErrNotFound = dao.ErrNotFound

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (rep *UserRepository) Create(ctx *gin.Context, user *domain.User) error {
	return rep.dao.Insert(ctx, dao.User{
		Password: user.Password,
		Email:    user.Email,
	})
}

func (rep *UserRepository) FindByEmail(ctx *gin.Context, user domain.User) (domain.User, error) {
	u, err := rep.dao.FindByEmail(ctx, dao.User{
		Email: user.Email,
	})
	uDomain := domain.User{
		Email:    u.Email,
		Password: u.Password,
		ID:       u.ID,
	}
	return uDomain, err
}
