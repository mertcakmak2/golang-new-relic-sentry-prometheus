package domain

import (
	"time"
)

type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Age         int       `json:"age"`
	CreatedDate time.Time `json:"created_date"`
}

type UserUseCase interface {
	CreateUser(user User) (User, *AppError)
	GetUserById(id uint) (User, *AppError)
	UpdateUser(user User) (User, *AppError)
	DeleteUserById(id uint) *AppError
}

type UserRepository interface {
	CreateUser(user User) (User, *AppError)
	GetUserById(id uint) (User, *AppError)
	UpdateUser(user User) (User, *AppError)
	DeleteUserById(id uint) *AppError
}
