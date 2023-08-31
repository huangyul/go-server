package dao

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx *gin.Context, u User) error {
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if mysqlErr.Number == 1062 {
			return errors.New("邮箱冲突")
		}
	}
	return err
}

type User struct {
	ID       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
}
