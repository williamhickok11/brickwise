package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError represents a structured API error response
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e *APIError) Error() string {
	return e.Message
}

// Predefined error types
var (
	ErrNotFound     = &APIError{Code: "NOT_FOUND", Message: "Resource not found", Status: http.StatusNotFound}
	ErrInvalidInput = &APIError{Code: "INVALID_INPUT", Message: "Invalid input provided", Status: http.StatusBadRequest}
	ErrConflict     = &APIError{Code: "CONFLICT", Message: "Resource conflict", Status: http.StatusConflict}
	ErrInternal     = &APIError{Code: "INTERNAL_ERROR", Message: "Internal server error", Status: http.StatusInternalServerError}
)

// New creates a new APIError with custom message
func New(code string, message string, status int) *APIError {
	return &APIError{Code: code, Message: message, Status: status}
}

// Wrap wraps an existing error with an APIError
func Wrap(err error, apiErr *APIError) *APIError {
	if err == nil {
		return nil
	}
	return &APIError{
		Code:    apiErr.Code,
		Message: fmt.Sprintf("%s: %v", apiErr.Message, err),
		Status:  apiErr.Status,
	}
}

// WriteJSON writes the error as JSON to the response writer
func WriteJSON(w http.ResponseWriter, err *APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err)
}
