/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package allerror provides a set of error codes and error types used in the application.
package allerror

import "strings"

const (
	errorCodeNoPermission = "no_permission"

	// ErrorCodeUserNotFound is const
	ErrorCodeUserNotFound = "user_not_found"

	// ErrorRateLimitOver is const
	ErrorRateLimitOver = "rate_limit_over"

	// ErrorCodeModelNotFound is const
	ErrorCodeModelNotFound = "model_not_found"

	// ErrorCodeSpaceNotFound is const
	ErrorCodeSpaceNotFound = "space_not_found"

	// ErrorCodeTokenNotFound is const
	ErrorCodeTokenNotFound = "token_not_found"

	// ErrorCodeRepoNotFound is const
	ErrorCodeRepoNotFound = "repo_not_found"

	// ErrorCodeOrganizationNotFound is const
	ErrorCodeOrganizationNotFound = "organization_not_found"

	// ErrorCodeCountExceeded is const
	ErrorCodeCountExceeded = "count_exceeded"

	// ErrorCodeSpaceAppNotFound space app
	ErrorCodeSpaceAppNotFound = "space_app_not_found"

	// ErrorCodeSpaceAppUnmatchedStatus is const
	ErrorCodeSpaceAppUnmatchedStatus = "space_app_unmatched_status"

	// ErrorCodeAccessTokenInvalid This error code is for restful api
	ErrorCodeAccessTokenInvalid = "access_token_invalid"

	// ErrorCodeCSRFTokenMissing is const
	ErrorCodeCSRFTokenMissing = "csrf_token_missing" // #nosec G101

	// ErrorCodeCSRFTokenInvalid is const
	ErrorCodeCSRFTokenInvalid = "csrf_token_invalid" // #nosec G101

	// ErrorCodeCSRFTokenNotFound is const
	ErrorCodeCSRFTokenNotFound = "csrf_token_not_found" // #nosec G101

	// ErrorCodeSessionInvalid is const
	ErrorCodeSessionInvalid = "session_invalid"

	// ErrorCodeSessionIdInvalid is const
	ErrorCodeSessionIdInvalid = "session_id_invalid"

	// ErrorCodeSessionIdMissing is const
	ErrorCodeSessionIdMissing = "session_id_missing"

	// ErrorCodeSessionNotFound is const
	ErrorCodeSessionNotFound = "session_not_found"

	// ErrorCodeBranchExist is const
	ErrorCodeBranchExist = "branch_exist"

	// ErrorCodeBranchInavtive is const
	ErrorCodeBranchInavtive = "branch_inactive"

	// ErrorCodeBaseBranchNotFound is const
	ErrorCodeBaseBranchNotFound = "base_branch_not_found"

	// ErrorCodeOrgExistResource is const
	ErrorCodeOrgExistResource = "org_resource_exist"

	errorCodeInvalidParam = "invalid_param"

	// ErrorEmailError is const
	ErrorEmailError = "email_error"

	// ErrorEmailCodeError is const
	ErrorEmailCodeError = "email_verify_code_error"

	// ErrorEmailCodeInvalid is const
	ErrorEmailCodeInvalid = "email_verify_code_invalid"

	// ErrorCodeNeedBindEmail is const
	ErrorCodeNeedBindEmail = "user_no_email"

	// ErrorCodeUserDuplicateBind is const
	ErrorCodeUserDuplicateBind = "user_duplicate_bind"

	// ErrorCodeEmailDuplicateBind is const
	ErrorCodeEmailDuplicateBind = "email_duplicate_bind"

	// ErrorCodeEmailDuplicateSend is const
	ErrorCodeEmailDuplicateSend = "email_duplicate_send"

	// ErrorBaseCase is const
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

// limitRateError
type limitRateError struct {
	errorImpl
}

// OverLimit is a marker method for over limit rate error.
func (l limitRateError) OverLimit() {}

// NewOverLimit creates a new over limit error with the specified code and message.
func NewOverLimit(code string, msg string) limitRateError {
	return limitRateError{errorImpl: New(code, msg)}
}
