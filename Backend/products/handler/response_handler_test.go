package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	customerrors "web-service/errors"

	"github.com/stretchr/testify/assert"
)

func TestHandleError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		message    string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "DatabaseError",
			err:        customerrors.NewDatabaseError("DB error", nil),
			message:    "Test message",
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"DB error","details":"Error: DB error, Cause: <nil>"}`,
		},
		{
			name:       "NotFoundError",
			err:        customerrors.NewNotFoundError("Not found", nil),
			message:    "Test message",
			wantStatus: http.StatusNotFound,
			wantBody:   `{"error":"Not found","details":"Error: Not found, Cause: <nil>"}`,
		},
		{
			name:       "S3Error",
			err:        customerrors.NewS3Error("S3 error", nil),
			message:    "Test message",
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"S3 error","details":"Error: S3 error, Cause: <nil>"}`,
		},
		{
			name:       "BadRequestError",
			err:        customerrors.NewBadRequestError("Bad request", nil),
			message:    "Test message",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"Bad request","details":"Error: Bad request, Cause: <nil>"}`,
		},
		{
			name:       "DefaultError",
			err:        fmt.Errorf("generic error"),
			message:    "Generic message",
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"Generic message","details":"generic error"}`,
		},
		{
			name:       "NilError",
			err:        nil,
			message:    "Nil message",
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"Nil message"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			HandleError(rr, tt.err, tt.message)

			assert.Equal(t, tt.wantStatus, rr.Code)
			assert.JSONEq(t, tt.wantBody, rr.Body.String())
		})
	}
}

func TestHandleSuccessResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		data       interface{}
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Success with data",
			statusCode: http.StatusOK,
			data:       map[string]string{"message": "success"},
			wantStatus: http.StatusOK,
			wantBody:   `{"message":"success"}`,
		},
		{
			name:       "Success with empty data",
			statusCode: http.StatusCreated,
			data:       nil,
			wantStatus: http.StatusCreated,
			wantBody:   `null`,
		},
		{
			name:       "Success with integer data",
			statusCode: http.StatusOK,
			data:       123,
			wantStatus: http.StatusOK,
			wantBody:   `123`,
		},
		{
			name:       "Success with string data",
			statusCode: http.StatusOK,
			data:       "test",
			wantStatus: http.StatusOK,
			wantBody:   `"test"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			HandleSuccessResponse(rr, tt.statusCode, tt.data)

			assert.Equal(t, tt.wantStatus, rr.Code)
			assert.JSONEq(t, tt.wantBody, rr.Body.String())
		})
	}
}

func TestHandleSuccessResponse_EncodingError(t *testing.T) {
	rr := httptest.NewRecorder()
	invalidData := make(chan int)

	HandleSuccessResponse(rr, http.StatusInternalServerError, invalidData)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.JSONEq(t, `{"error":"Error encoding response","details":"json: unsupported type: chan int"}`, rr.Body.String())
}
