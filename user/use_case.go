package user

import (
	"fmt"
	"go-app/domain"
	"go.uber.org/zap"
	"time"
)

//go:generate mockgen -destination=../mocks/mockUserUsecase.go -package=mocks go-app/domain UserUseCase
type userUseCase struct {
	repo   domain.UserRepository
	logger *zap.Logger
}

func NewUserUseCase(repo domain.UserRepository, logger *zap.Logger) domain.UserUseCase {
	return &userUseCase{repo: repo, logger: logger}
}

func (u *userUseCase) CreateUser(user domain.User) (domain.User, *domain.AppError) {
	user.CreatedDate = time.Now()
	if user.Name == "" {
		err := domain.NewValidationError("The name should not be empty.")
		u.logger.Error(err.Message)
		return user, err
	}

	createdUser, err := u.repo.CreateUser(user)
	if err != nil {
		u.logger.Error(err.Message)
		return domain.User{}, err
	}

	u.logger.Info(fmt.Sprintf("User created. ID: %d", createdUser.ID))
	return createdUser, nil
}

func (u *userUseCase) GetUserById(id uint) (domain.User, *domain.AppError) {
	user, err := u.repo.GetUserById(id)
	if err != nil {
		u.logger.Error(err.Message)
		return user, err
	}

	return user, nil
}

func (u *userUseCase) UpdateUser(user domain.User) (domain.User, *domain.AppError) {
	updatedUser, err := u.repo.UpdateUser(user)
	if err != nil {
		u.logger.Error(err.Message)
		return updatedUser, err
	}
	return updatedUser, nil
}

func (u *userUseCase) DeleteUserById(id uint) *domain.AppError {
	err := u.repo.DeleteUserById(id)
	if err != nil {
		u.logger.Error(err.Message)
		return err
	}
	return err
}
