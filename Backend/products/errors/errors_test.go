package customerrors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomError(t *testing.T) {
	cause := errors.New("test cause")
	err := NewCustomError("test message", http.StatusInternalServerError, cause)

	assert.Equal(t, "test message", err.GetMessage()) // Use getter method
	assert.Equal(t, http.StatusInternalServerError, err.GetStatusCode())
	assert.Equal(t, cause, err.GetCause())
	assert.Equal(t, "Error: test message, Cause: test cause", err.Error())

	errNilCause := NewCustomError("test message", http.StatusInternalServerError, nil)
	assert.Equal(t, "Error: test message, Cause: <nil>", errNilCause.Error())
}

func TestNotFoundError(t *testing.T) {
	cause := errors.New("test cause")
	err := NewNotFoundError("not found", cause)

	assert.Equal(t, "not found", err.GetMessage()) // Use getter
	assert.Equal(t, http.StatusNotFound, err.GetStatusCode())
	assert.Equal(t, cause, err.GetCause())
	assert.Equal(t, "Error: not found, Cause: test cause", err.Error())
}

func TestDatabaseError(t *testing.T) {
	cause := errors.New("test cause")
	err := NewDatabaseError("database error", cause)

	assert.Equal(t, "database error", err.GetMessage())
	assert.Equal(t, http.StatusInternalServerError, err.GetStatusCode())
	assert.Equal(t, cause, err.GetCause())
	assert.Equal(t, "Error: database error, Cause: test cause", err.Error())
}

func TestS3Error(t *testing.T) {
	cause := errors.New("test cause")
	err := NewS3Error("s3 error", cause)

	assert.Equal(t, "s3 error", err.GetMessage())
	assert.Equal(t, http.StatusInternalServerError, err.GetStatusCode())
	assert.Equal(t, cause, err.GetCause())
	assert.Equal(t, "Error: s3 error, Cause: test cause", err.Error())
}

func TestBadRequestError(t *testing.T) {
	cause := errors.New("test cause")
	err := NewBadRequestError("bad request", cause)

	assert.Equal(t, "bad request", err.GetMessage())
	assert.Equal(t, http.StatusBadRequest, err.GetStatusCode())
	assert.Equal(t, cause, err.GetCause())
	assert.Equal(t, "Error: bad request, Cause: test cause", err.Error())
}

func TestErrorWithNilCause(t *testing.T) {
	err := NewBadRequestError("bad request", nil)

	assert.Equal(t, "Error: bad request, Cause: <nil>", err.Error())
}
