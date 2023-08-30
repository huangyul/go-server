package dao

import "gorm.io/gorm"

func InitUser(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
