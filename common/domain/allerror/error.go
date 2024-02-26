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

	// space app
	ErrorCodeSpaceAppNotFound        = "space_app_not_found"
	ErrorCodeSpaceAppUnmatchedStatus = "space_app_unmatched_status"

	// This error code is for restful api
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
	return notfoudError{errorImpl: New(code, msg)}
}

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}

	_, ok := err.(notfoudError)

	return ok
}

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

func (e noPermissionError) NoPermission() {}

// NewNoPermission
func NewNoPermission(msg string) noPermissionError {
	return noPermissionError{errorImpl: New(errorCodeNoPermission, msg)}
}

func IsNoPermission(err error) bool {
	if err == nil {
		return false
	}

	_, ok := err.(noPermissionError)

	return ok
}

func NewInvalidParam(msg string) errorImpl {
	return New(errorCodeInvalidParam, msg)
}

func NewCountExceeded(msg string) errorImpl {
	return New(ErrorCodeCountExceeded, msg)
}
