package user

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go-app/config"
	"go-app/domain"
	"go-app/mocks"
	"testing"
)

var (
	_userMockRepo *mocks.MockUserRepository
	_userUseCase  domain.UserUseCase
)

func mockUseCaseSetup(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	// Mock UserRepository
	_userMockRepo = mocks.NewMockUserRepository(c)

	logger := config.ZapTestConfig()
	_userUseCase = NewUserUseCase(_userMockRepo, logger)
}

func Test_Should_Create_User_With_MockUserRepository(t *testing.T) {
	mockUseCaseSetup(t)

	// GIVEN
	user := domain.User{Name: "test", Age: 18}
	expectedUser := domain.User{ID: 1, Name: "test", Age: 18}

	// WHEN
	_userMockRepo.EXPECT().CreateUser(gomock.Any()).Return(expectedUser, nil)
	res, err := _userUseCase.CreateUser(user)

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, expectedUser.Name, res.Name)
}

func Test_Should_Return_Validation_Err_When_Invoke_Create_User_With_MockUserRepository(t *testing.T) {
	mockUseCaseSetup(t)

	// GIVEN
	user := domain.User{Age: 18}
	validationErr := domain.NewValidationError("The name should not be empty.")

	// WHEN
	_, err := _userUseCase.CreateUser(user)

	// THEN
	assert.NotNil(t, err)
	assert.Equal(t, validationErr.Message, err.Message)
}

func Test(t *testing.T) {
	mockUseCaseSetup(t)

	// GIVEN
	user := domain.User{Name: "test-user", Age: 18}
	expectedErr := domain.NewUnexpectedError("Unexpected error.")

	// WHEN
	_userMockRepo.EXPECT().CreateUser(gomock.Any()).Return(domain.User{}, expectedErr)
	_, err := _userUseCase.CreateUser(user)

	// THEN
	assert.NotNil(t, err)
	assert.Equal(t, expectedErr.Message, err.Message)
}

func Test_Should_Return_Unexpected_Err_When_Invoke_Create_User_With_MockUserRepository(t *testing.T) {
	mockUseCaseSetup(t)

	// GIVEN
	expectedUser := domain.User{ID: 1, Name: "test", Age: 18}
	var id uint = 1

	// WHEN
	_userMockRepo.EXPECT().GetUserById(id).Return(expectedUser, nil)
	res, err := _userUseCase.GetUserById(id)

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, expectedUser.Name, res.Name)
}

func Test_Should_Return_Not_Found_Err_When_Invoke_Get_User_By_Id_With_MockUserRepository(t *testing.T) {
	mockUseCaseSetup(t)

	// GIVEN
	var id uint = 1
	errStr := fmt.Sprintf("User not found, ID: %d", id)
	notFoundErr := domain.NewNotFoundError(errStr)

	// WHEN
	_userMockRepo.EXPECT().GetUserById(gomock.Any()).Return(domain.User{}, notFoundErr)
	_, err := _userUseCase.GetUserById(id)

	// THEN
	assert.NotNil(t, err)
	assert.Equal(t, notFoundErr.Message, err.Message)
}

func Test_Should_Update_User_With_MockUserRepository(t *testing.T) {
	mockUseCaseSetup(t)

	// GIVEN
	user := domain.User{ID: 1, Name: "updated-user", Age: 18}
	expectedUser := domain.User{ID: 1, Name: "updated-user", Age: 18}

	// WHEN
	_userMockRepo.EXPECT().UpdateUser(gomock.Any()).Return(expectedUser, nil)
	res, err := _userUseCase.UpdateUser(user)

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, expectedUser.ID, res.ID)
	assert.Equal(t, expectedUser.Name, res.Name)
}

func Test_Should_Return_Unexpected_Err_When_Invoke_Update_User_With_MockUserRepository(t *testing.T) {
	mockUseCaseSetup(t)

	// GIVEN
	user := domain.User{ID: 1, Name: "updated-user", Age: 18}
	errStr := fmt.Sprintf("Unexpected Error")
	expectedErr := domain.NewUnexpectedError(errStr)

	// WHEN
	_userMockRepo.EXPECT().UpdateUser(gomock.Any()).Return(domain.User{}, expectedErr)
	_, err := _userUseCase.UpdateUser(user)

	// THEN
	assert.NotNil(t, err)
	assert.Equal(t, expectedErr.Message, err.Message)
}

func Test_Should_Delete_User_By_Id_With_MockUserRepository(t *testing.T) {
	mockUseCaseSetup(t)

	// GIVEN
	var id uint = 1

	// WHEN
	_userMockRepo.EXPECT().DeleteUserById(gomock.Any()).Return(nil)
	err := _userUseCase.DeleteUserById(id)

	// THEN
	assert.Nil(t, err)
}

func Test_Should_Return_Unexpected_Err_When_Invoke_Delete_User_By_Id_With_MockUserRepository(t *testing.T) {
	mockUseCaseSetup(t)

	// GIVEN
	var id uint = 1
	errStr := fmt.Sprintf("Unexpected Error")
	expectedErr := domain.NewUnexpectedError(errStr)

	// WHEN
	_userMockRepo.EXPECT().DeleteUserById(gomock.Any()).Return(expectedErr)
	err := _userUseCase.DeleteUserById(id)

	// THEN
	assert.NotNil(t, err)
	assert.Equal(t, expectedErr.Message, err.Message)
}
