package repository

// ErrorDuplicateCreating
type ErrorDuplicateCreating struct {
	error
}

func NewErrorDuplicateCreating(err error) ErrorDuplicateCreating {
	return ErrorDuplicateCreating{error: err}
}

// ErrorResourceNotExists
type ErrorResourceNotExists struct {
	error
}

func NewErrorResourceNotExists(err error) ErrorResourceNotExists {
	return ErrorResourceNotExists{error: err}
}

// ErrorConcurrentUpdating
type ErrorConcurrentUpdating struct {
	error
}

func NewErrorConcurrentUpdating(err error) ErrorConcurrentUpdating {
	return ErrorConcurrentUpdating{error: err}
}

// helper

func IsErrorResourceNotExists(err error) bool {
	_, ok := err.(ErrorResourceNotExists)

	return ok
}

func IsErrorDuplicateCreating(err error) bool {
	_, ok := err.(ErrorDuplicateCreating)

	return ok
}

func IsErrorConcurrentUpdating(err error) bool {
	_, ok := err.(ErrorConcurrentUpdating)

	return ok
}
