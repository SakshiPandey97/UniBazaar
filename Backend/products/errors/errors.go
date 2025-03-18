package customerrors

import (
	"fmt"
	"net/http"
)

type CustomErrorInterface interface {
	Error() string
	GetStatusCode() int
	GetMessage() string
	GetCause() error
}

type CustomError struct {
	Message    string
	StatusCode int
	Cause      error
}

func NewCustomError(message string, statusCode int, cause error) *CustomError {
	return &CustomError{
		Message:    message,
		StatusCode: statusCode,
		Cause:      cause,
	}
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("Error: %s, Cause: %v", e.Message, e.Cause)
}

func (e *CustomError) GetStatusCode() int {
	return e.StatusCode
}

func (e *CustomError) GetMessage() string {
	return e.Message
}

func (e *CustomError) GetCause() error {
	return e.Cause
}

type NotFoundError struct {
	*CustomError
}

func NewNotFoundError(message string, cause error) *NotFoundError {
	return &NotFoundError{
		&CustomError{
			Message:    message,
			StatusCode: http.StatusNotFound,
			Cause:      cause,
		},
	}
}

type DatabaseError struct {
	*CustomError
}

func NewDatabaseError(message string, cause error) *DatabaseError {
	return &DatabaseError{
		&CustomError{
			Message:    message,
			StatusCode: http.StatusInternalServerError,
			Cause:      cause,
		},
	}
}

type S3Error struct {
	*CustomError
}

func NewS3Error(message string, cause error) *S3Error {
	return &S3Error{
		&CustomError{
			Message:    message,
			StatusCode: http.StatusInternalServerError,
			Cause:      cause,
		},
	}
}

type BadRequestError struct {
	*CustomError
}

func NewBadRequestError(message string, cause error) *BadRequestError {
	return &BadRequestError{
		&CustomError{
			Message:    message,
			StatusCode: http.StatusBadRequest,
			Cause:      cause,
		},
	}
}
