package utils

import "fmt"

// NotFoundError represents an error when a requested resource is not found
type NotFoundError struct {
	Resource string
	ID       interface{}
}

// Error implements the error interface
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found: %v", e.Resource, e.ID)
}

// IsNotFoundError checks if an error is a NotFoundError
func IsNotFoundError(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}

// NewNotFoundError creates a new NotFoundError
func NewNotFoundError(resource string, id interface{}) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
	}
}