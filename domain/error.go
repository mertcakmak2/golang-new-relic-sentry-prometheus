package domain

import (
	"net/http"
)

type AppError struct {
	Code    int    `json:",omitempty"`
	Message string `json:"message"`
}

func (e AppError) AsMessageError() *AppError {
	return &AppError{
		Message: e.Message,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: message}
}

func NewUnexpectedError(message string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: message}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: message}
}

func NewUserAlreadyExistError(message string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: message}
}

func NewValidationError(message string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: message}
}
