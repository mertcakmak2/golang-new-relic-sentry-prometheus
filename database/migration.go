package database

import (
	"go-app/domain"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&domain.User{})
}
