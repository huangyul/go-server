package dao

import "gorm.io/gorm"

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

type User struct {
	ID int64 `gorm:"primaryKey,autoIncrement"`
}
