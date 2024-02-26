/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repository provides custom error types and utility functions for error handling.
package repository

// ErrorDuplicateCreating represents an error indicating a duplicate creation attempt.
type ErrorDuplicateCreating struct {
	error
}

// NewErrorDuplicateCreating creates a new ErrorDuplicateCreating error with the given underlying error.
func NewErrorDuplicateCreating(err error) ErrorDuplicateCreating {
	return ErrorDuplicateCreating{error: err}
}

// ErrorResourceNotExists represents an error indicating a non-existent resource.
type ErrorResourceNotExists struct {
	error
}

// NewErrorResourceNotExists creates a new ErrorResourceNotExists error with the given underlying error.
func NewErrorResourceNotExists(err error) ErrorResourceNotExists {
	return ErrorResourceNotExists{error: err}
}

// ErrorConcurrentUpdating represents an error indicating a concurrent update conflict.
type ErrorConcurrentUpdating struct {
	error
}

// NewErrorConcurrentUpdating creates a new ErrorConcurrentUpdating error with the given underlying error.
func NewErrorConcurrentUpdating(err error) ErrorConcurrentUpdating {
	return ErrorConcurrentUpdating{error: err}
}

// IsErrorResourceNotExists checks if the given error is of type ErrorResourceNotExists.
func IsErrorResourceNotExists(err error) bool {
	_, ok := err.(ErrorResourceNotExists)

	return ok
}

// IsErrorDuplicateCreating checks if the given error is of type ErrorDuplicateCreating.
func IsErrorDuplicateCreating(err error) bool {
	_, ok := err.(ErrorDuplicateCreating)

	return ok
}

// IsErrorConcurrentUpdating checks if the given error is of type ErrorConcurrentUpdating.
func IsErrorConcurrentUpdating(err error) bool {
	_, ok := err.(ErrorConcurrentUpdating)

	return ok
}
