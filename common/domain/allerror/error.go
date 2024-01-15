package allerror

import "strings"

const (
	errorCodeNoPermission = "no_permission"

	ErrorCodeModelNotFound = "model_not_found"

	ErrorCodeSpaceNotFound = "space_not_found"

	ErrorCodeUserNotFound = "user_not_found"

	ErrorCodeTokenNotFound = "token_not_found"

	ErrorCodeOrganizationNotFound = "organization_not_found"

	// This error code is for restful api
	ErrorCodeAccessTokenInvalid = "access_token_invalid"

	ErrorCodeLoginIdInvalid    = "login_id_invalid"
	ErrorCodeLoginIdMissing    = "login_id_missing"
	ErrorCodeLoginIdNotFound   = "login_id_not_found"
	ErrorCodeCSRFTokenMissing  = "csrf_token_missing"   // #nosec G101
	ErrorCodeCSRFTokenInvalid  = "csrf_token_invalid"   // #nosec G101
	ErrorCodeCSRFTokenNotFound = "csrf_token_not_found" // #nosec G101

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
