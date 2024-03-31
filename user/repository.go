package user

import (
	"errors"
	"fmt"
	"go-app/domain"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../mocks/mockUserRepository.go -package=mocks go-app/domain UserRepository
type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user domain.User) (domain.User, *domain.AppError) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, domain.NewUnexpectedError(err.Error())
	}
	return user, nil
}

func (r *userRepository) GetUserById(id uint) (domain.User, *domain.AppError) {
	var user domain.User
	// err := r.db.First(&user, id).Error
	err := r.db.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		errStr := fmt.Sprintf("User not found, ID: %d", id)
		return user, domain.NewNotFoundError(errStr)
	}

	if err != nil {
		return user, domain.NewUnexpectedError(err.Error())
	}

	return user, nil
}

func (r *userRepository) UpdateUser(user domain.User) (domain.User, *domain.AppError) {
	// err := r.db.WithContext(context.Background()).Model(user).Where("id = ?", user.ID).Update("name", user.Name).Error
	err := r.db.Save(&user).Error
	if err != nil {
		return user, domain.NewUnexpectedError(err.Error())
	}
	return user, nil
}

func (r *userRepository) DeleteUserById(id uint) *domain.AppError {
	err := r.db.Delete(&domain.User{}, id).Error
	if err != nil {
		return domain.NewUnexpectedError(err.Error())
	}
	return nil
}
