package repositories

// errorDuplicateCreating
type errorDuplicateCreating struct {
	error
}

func NewErrorDuplicateCreating(err error) errorDuplicateCreating {
	return errorDuplicateCreating{err}
}

// errorDataNotExists
type errorDataNotExists struct {
	error
}

func NewErrorDataNotExists(err error) errorDataNotExists {
	return errorDataNotExists{err}
}

func isErrorDataNotExists(err error) bool {
	_, ok := err.(errorDataNotExists)

	return ok
}

// errorConcurrentUpdating
type errorConcurrentUpdating struct {
	error
}

func NewErrorConcurrentUpdating(err error) errorConcurrentUpdating {
	return errorConcurrentUpdating{err}
}
