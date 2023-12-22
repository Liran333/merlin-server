package allerror

import "strings"

const (
	errorCodeNoPermission = "no_permission"

	ErrorCodeModelNotFound = "model_not_found"

	ErrorCodeUserNotFound = "user_not_found"

	ErrorCodeOrganizationNotFound = "organization_not_found"

	ErrorCodeAccessTokenInvalid = "access_token_invalid"

	errorCodeInvalidParam = "invalid_param"
)

// errorImpl
type errorImpl struct {
	code string
	msg  string
}

func (e errorImpl) Error() string {
	return e.msg
}

func (e errorImpl) ErrorCode() string {
	return e.code
}

// New
func New(code string, msg string) errorImpl {
	v := errorImpl{
		code: code,
	}

	if msg == "" {
		v.msg = strings.ReplaceAll(code, "_", " ")
	} else {
		v.msg = msg
	}

	return v
}

// notfoudError
type notfoudError struct {
	errorImpl
}

func (e notfoudError) NotFound() {}

// NewNotFound
func NewNotFound(code string, msg string) notfoudError {
	return notfoudError{New(code, msg)}
}

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}

	if _, ok := err.(notfoudError); ok {
		return true
	}

	return false
}

// noPermissionError
type noPermissionError struct {
	errorImpl
}

func (e noPermissionError) NoPermission() {}

// NewNoPermission
func NewNoPermission(msg string) noPermissionError {
	return noPermissionError{New(errorCodeNoPermission, msg)}
}

func NewInvalidParam(msg string) errorImpl {
	return New(errorCodeInvalidParam, msg)
}
