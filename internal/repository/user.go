package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-server/internal/domain"
	"github.com/huangyul/go-server/internal/repository/dao"
)

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
