package errors

import (
	stderrors "errors"
	"fmt"
	"net/http"
	"testing"
)

func TestHttpError_Error(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
		details    map[string]interface{}
		err        error
		want       string
	}{
		{
			name:       "Basic error without original error",
			statusCode: 404,
			message:    "Resource not found",
			details:    nil,
			err:        nil,
			want:       "HTTP Error 404: Resource not found",
		},
		{
			name:       "Error with original error",
			statusCode: 500,
			message:    "Internal server error",
			details:    nil,
			err:        stderrors.New("database connection failed"),
			want:       "HTTP Error 500: Internal server error - database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &HttpError{
				StatusCode: tt.statusCode,
				Message:    tt.message,
				Details:    tt.details,
				Err:        tt.err,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("HttpError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewHttpError(t *testing.T) {
	statusCode := 400
	message := "Bad request"
	details := map[string]interface{}{
		"field": "email",
		"error": "invalid format",
	}
	originalErr := stderrors.New("validation failed")

	err := NewHttpError(statusCode, message, details, originalErr)

	if err.StatusCode != statusCode {
		t.Errorf("NewHttpError().StatusCode = %v, want %v", err.StatusCode, statusCode)
	}
	if err.Message != message {
		t.Errorf("NewHttpError().Message = %v, want %v", err.Message, message)
	}
	if fmt.Sprintf("%v", err.Details) != fmt.Sprintf("%v", details) {
		t.Errorf("NewHttpError().Details = %v, want %v", err.Details, details)
	}
	if err.Err != originalErr {
		t.Errorf("NewHttpError().Err = %v, want %v", err.Err, originalErr)
	}
}

func TestBadRequest(t *testing.T) {
	message := "Bad request"
	err := BadRequest(message)

	if err.StatusCode != http.StatusBadRequest {
		t.Errorf("BadRequest().StatusCode = %v, want %v", err.StatusCode, http.StatusBadRequest)
	}
	if err.Message != message {
		t.Errorf("BadRequest().Message = %v, want %v", err.Message, message)
	}
}
