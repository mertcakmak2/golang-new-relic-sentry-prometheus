package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go-app/config"
	"go-app/domain"
	"go-app/mocks"
	"go-app/user"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	_userMockUseCase *mocks.MockUserUseCase
	_userHandler     *user.Handler
)

func handlerSetupRouter(t *testing.T) *gin.Engine {
	c := gomock.NewController(t)
	defer c.Finish()

	// Mock UserUseCase
	_userMockUseCase = mocks.NewMockUserUseCase(c)

	logger := config.ZapTestConfig()
	_userHandler = user.NewUserHandler(_userMockUseCase, logger)

	r := setupRouter(config.NewRelicConfig(), _userHandler)
	return r

}

func Test_Should_Create_User_With_MockUserUseCase(t *testing.T) {
	router := handlerSetupRouter(t)

	// GIVEN
	u := domain.User{Name: "created-user", Age: 22}
	byteUser, _ := json.Marshal(u)
	expectedUser := domain.User{ID: 10, Name: u.Name, Age: u.Age}

	// WHEN
	_userMockUseCase.EXPECT().CreateUser(gomock.Any()).Return(expectedUser, nil)

	w := httptest.NewRecorder()
	url := "/api/v1/users"
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(byteUser))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// THEN
	savedUser := domain.User{}

	assert.NotEmpty(t, w.Body.String())
	err := json.Unmarshal([]byte(w.Body.String()), &savedUser)

	assert.Nil(t, err)
	assert.Equal(t, 201, w.Code)
	assert.Equal(t, expectedUser.Name, savedUser.Name)
	assert.NotEmpty(t, savedUser.ID)
}

func Test_Should_Return_Unexpected_Err_When_Invoke_Create_User_With_MockUserUseCase(t *testing.T) {
	router := handlerSetupRouter(t)

	// GIVEN
	u := domain.User{Name: "created-user", Age: 22}
	byteUser, _ := json.Marshal(u)

	gormErr := errors.New("Unexpected Error")
	expectedErr := domain.NewUnexpectedError(gormErr.Error())

	// WHEN
	_userMockUseCase.EXPECT().CreateUser(gomock.Any()).Return(domain.User{}, expectedErr)

	w := httptest.NewRecorder()
	url := "/api/v1/users"
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(byteUser))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// THEN
	resErr := domain.AppError{}

	assert.NotEmpty(t, w.Body.String())
	err := json.Unmarshal([]byte(w.Body.String()), &resErr)

	assert.Nil(t, err)
	assert.Equal(t, 500, w.Code)
	assert.Equal(t, expectedErr.Message, resErr.Message)
}

func Test_Should_Find_User_With_MockUserUseCase(t *testing.T) {
	router := handlerSetupRouter(t)

	// GIVEN
	var id uint = 1
	expectedUser := domain.User{ID: id, Name: "test", Age: 18}

	// WHEN
	_userMockUseCase.EXPECT().GetUserById(gomock.Any()).Return(expectedUser, nil)

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/api/v1/users/%d", id)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	router.ServeHTTP(w, req)

	// THEN
	u := domain.User{}

	assert.NotEmpty(t, w.Body.String())
	err := json.Unmarshal([]byte(w.Body.String()), &u)

	assert.Nil(t, err)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, id, u.ID)
}

func Test_Should_Return_Not_Found_Err_When_Invoke_Find_User_With_MockUserUseCase(t *testing.T) {
	router := handlerSetupRouter(t)

	// GIVEN
	var id uint = 1
	errStr := fmt.Sprintf("User not found, ID: %d", id)
	expectedErr := domain.NewNotFoundError(errStr)

	// WHEN
	_userMockUseCase.EXPECT().GetUserById(gomock.Any()).Return(domain.User{}, expectedErr)

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/api/v1/users/%d", id)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	router.ServeHTTP(w, req)

	// THEN
	resErr := domain.AppError{}

	assert.NotEmpty(t, w.Body.String())
	_ = json.Unmarshal([]byte(w.Body.String()), &resErr)

	assert.NotNil(t, resErr)
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, expectedErr.Message, resErr.Message)
}

func Test_Should_Update_User_With_MockUserUseCase(t *testing.T) {
	router := handlerSetupRouter(t)

	// GIVEN
	expectedUser := domain.User{ID: 5, Name: "updated-user", Age: 22}
	byteUser, _ := json.Marshal(expectedUser)

	// WHEN
	_userMockUseCase.EXPECT().UpdateUser(gomock.Any()).Return(expectedUser, nil)

	w := httptest.NewRecorder()
	url := "/api/v1/users"
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(byteUser))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// THEN
	updatedUser := domain.User{}

	assert.NotEmpty(t, w.Body.String())
	err := json.Unmarshal([]byte(w.Body.String()), &updatedUser)

	assert.Nil(t, err)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, expectedUser.Name, updatedUser.Name)
	assert.NotEmpty(t, updatedUser.ID)
}

func Test_Should_Return_Unexpected_Err_When_Invoke_Update_User_With_MockUserUseCase(t *testing.T) {
	router := handlerSetupRouter(t)

	// GIVEN
	expectedUser := domain.User{ID: 5, Name: "updated-user", Age: 22}
	byteUser, _ := json.Marshal(expectedUser)

	gormErr := errors.New("Unexpected Error")
	expectedErr := domain.NewUnexpectedError(gormErr.Error())

	// WHEN
	_userMockUseCase.EXPECT().UpdateUser(gomock.Any()).Return(domain.User{}, expectedErr)

	w := httptest.NewRecorder()
	url := "/api/v1/users"
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(byteUser))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// THEN
	resErr := domain.AppError{}

	assert.NotEmpty(t, w.Body.String())
	err := json.Unmarshal([]byte(w.Body.String()), &resErr)

	assert.Nil(t, err)
	assert.Equal(t, 500, w.Code)
	assert.Equal(t, expectedErr.Message, resErr.Message)
}

func Test_Should_Delete_User_With_MockUserUseCase(t *testing.T) {
	router := handlerSetupRouter(t)

	// GIVEN
	var id uint = 1

	// WHEN
	_userMockUseCase.EXPECT().DeleteUserById(gomock.Any()).Return(nil)

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/api/v1/users/%d", id)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	router.ServeHTTP(w, req)

	// THEN
	assert.Equal(t, 204, w.Code)
}

func Test_Should_Return_Unexpected_Err_When_Invoke_Delete_User_With_MockUserUseCase(t *testing.T) {
	router := handlerSetupRouter(t)

	// GIVEN
	var id uint = 1
	gormErr := errors.New("Unexpected Error")
	expectedErr := domain.NewUnexpectedError(gormErr.Error())

	// WHEN
	_userMockUseCase.EXPECT().DeleteUserById(gomock.Any()).Return(expectedErr)

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/api/v1/users/%d", id)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	router.ServeHTTP(w, req)

	// THEN
	resErr := domain.AppError{}

	assert.NotEmpty(t, w.Body.String())
	err := json.Unmarshal([]byte(w.Body.String()), &resErr)

	assert.Nil(t, err)
	assert.Equal(t, 500, w.Code)
	assert.Equal(t, expectedErr.Message, resErr.Message)
}
