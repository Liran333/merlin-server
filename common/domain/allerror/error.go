/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package allerror provides a set of error codes and error types used in the application.
package allerror

import "strings"

const (
	errorCodeNoPermission = "no_permission"

	ErrorCodeUserNotFound = "user_not_found"

	ErrorCodeModelNotFound = "model_not_found"

	ErrorCodeSpaceNotFound = "space_not_found"

	ErrorCodeTokenNotFound = "token_not_found"

	ErrorCodeRepoNotFound = "repo_not_found"

	ErrorCodeOrganizationNotFound = "organization_not_found"

	ErrorCodeCountExceeded = "count_exceeded"

	// ErrorCodeSpaceAppNotFound space app
	ErrorCodeSpaceAppNotFound        = "space_app_not_found"
	ErrorCodeSpaceAppUnmatchedStatus = "space_app_unmatched_status"

	// ErrorCodeAccessTokenInvalid This error code is for restful api
	ErrorCodeAccessTokenInvalid = "access_token_invalid"

	ErrorCodeLoginIdInvalid    = "login_id_invalid"
	ErrorCodeLoginIdMissing    = "login_id_missing"
	ErrorCodeLoginIdNotFound   = "login_id_not_found"
	ErrorCodeCSRFTokenMissing  = "csrf_token_missing"   // #nosec G101
	ErrorCodeCSRFTokenInvalid  = "csrf_token_invalid"   // #nosec G101
	ErrorCodeCSRFTokenNotFound = "csrf_token_not_found" // #nosec G101

	ErrorCodeBranchExist        = "branch_exist"
	ErrorCodeBranchInavtive     = "branch_inactive"
	ErrorCodeBaseBranchNotFound = "base_branch_not_found"

	ErrorCodeOrgExistModel = "org_model_exist"

	errorCodeInvalidParam = "invalid_param"

	ErrorEmailError             = "email_error"
	ErrorEmailCodeError         = "email_verify_code_error"
	ErrorEmailCodeInvalid       = "email_verify_code_invalid"
	ErrorCodeNeedBindEmail      = "user_no_email"
	ErrorCodeUserDuplicateBind  = "user_duplicate_bind"
	ErrorCodeEmailDuplicateBind = "email_duplicate_bind"
	ErrorCodeEmailDuplicateSend = "email_duplicate_send"

	ErrorBaseCase = "internal_error"
)

// errorImpl
type errorImpl struct {
	code string
	msg  string
}

// Error returns the error message.
func (e errorImpl) Error() string {
	return e.msg
}

// ErrorCode returns the error code.
func (e errorImpl) ErrorCode() string {
	return e.code
}

// New creates a new error with the specified code and message.
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

// NotFound is a marker method for a not found error.
func (e notfoudError) NotFound() {}

// NewNotFound creates a new not found error with the specified code and message.
func NewNotFound(code string, msg string) notfoudError {
	return notfoudError{errorImpl: New(code, msg)}
}

// IsNotFound checks if the given error is a not found error.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}

	_, ok := err.(notfoudError)

	return ok
}

// IsUserDuplicateBind checks if the given error is a user duplicate bind error.
func IsUserDuplicateBind(err error) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(errorImpl); ok {
		if e.ErrorCode() == ErrorCodeUserDuplicateBind {
			return true
		}
	}

	return false
}

// noPermissionError
type noPermissionError struct {
	errorImpl
}

// NoPermission is a marker method for a "no permission" error.
func (e noPermissionError) NoPermission() {}

// NewNoPermission creates a new "no permission" error with the specified message.
func NewNoPermission(msg string) noPermissionError {
	return noPermissionError{errorImpl: New(errorCodeNoPermission, msg)}
}

// IsNoPermission checks if the given error is a "no permission" error.
func IsNoPermission(err error) bool {
	if err == nil {
		return false
	}

	_, ok := err.(noPermissionError)

	return ok
}

// NewInvalidParam creates a new error with the specified invalid parameter message.
func NewInvalidParam(msg string) errorImpl {
	return New(errorCodeInvalidParam, msg)
}

// NewCountExceeded creates a new error with the specified count exceeded message.
func NewCountExceeded(msg string) errorImpl {
	return New(ErrorCodeCountExceeded, msg)
}
