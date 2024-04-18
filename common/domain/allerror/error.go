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

	// ErrorCodeSpaceVariableNotFound space variable
	ErrorCodeSpaceVariableNotFound = "space_variable_not_found"

	// ErrorCodeSpaceSecretNotFound space secret
	ErrorCodeSpaceSecretNotFound = "space_secret_not_found"

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

	// ErrorCodeSpaceAppRestartOverTime is const
	ErrorCodeSpaceAppRestartOverTime = "space_app_restart_over_time"

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

	// ErrorCodeDisAgreedPrivacy is const
	ErrorCodeDisAgreedPrivacy = "disagreed_privacy"

	// ErrorCodeExpired
	ErrorCodeExpired = "expired"

	// privilege org id mismatch
	ErrorCodePrivilegeOrgIdMismatch = "privilege_org_id_mismatch"

	// not in a privilege org
	ErrorCodeNotInPrivilegeOrg = "not_in_privilege_org"

	// user has insufficient quota balance
	ErrorCodeInsufficientQuota = "insufficient_quota"

	// user is not currently using any quota
	ErrorCodeNoUsedQuota = "no_used_quota"

	// user has no npu permission
	ErrorCodeNoNpuPermission = "no_npu_permission"

	// ErrorBaseCase is const
	ErrorBaseCase = "internal_error"

	// ErrorMsgPublishFailed
	ErrorMsgPublishFailed = "msg_publish_failed"
	
	// dulicate creating
	ErrorDuplicateCreating = "duplicate_creating"

	// failed to get owner info
	ErrorFailedGetOwnerInfo = "failed_to_get_owner_info"

	// failed to get platform user
	ErrorFailGetPlatformUser = "failed_to_get_platform_user"

	// failed to create org
	ErrorFailedCreateOrg = "failed_to_create_org"

	// failed to create to org
	ErrorFailedCreateToOrg = "failed_to_create_to_org"

	// failed to save org member
	ErrorFailSaveOrgMember = "failed_to_save_org_member"

	// system error
	ErrorSystemError = "system_error"

	// when list member requests, org_name/requester can't be all empty
	ErrorOrgNameRequesterAllEmpty = "org_name_requester_all_empty"

	// only one of the org_name/invitee/inviter can be used
	ErrorOverOrgnameInviteeInviter = "over_orgname_invitee_inviter"

	// when list member invitation, org_name/invitee/inviter can't be all empty
	ErrorMemberInvitationParamAllEmpty = "member_invitation_param_all_empty"

	// failed to delete user in db
	ErrorFailedToDeleteUser = "failed_to_delete_user"

	// failed to get user info
	ErrorFailedToGetUserInfo = "failed_to_get_user_info"

	// name %s is already been taken
	ErrorNameAlreadyBeenTaken = "name_is_already_been_taken"

	// %s can't delete the org
	ErrorAccountCannotDeleteTheOrg = "account_can_not_delete_the_org"

	// failed to get org when get org by user, %w
	ErrorFailedToGetOrg = "failed_to_get_org"

	// missing name when creating token
	ErrorMissingName = "missing_name"

	// missing account when creating token
	ErrorMissingAccount = "missing_account"

	// failed to remove member
	ErrorFailedToRemoveMember = "failed_to_remove_member"

	// the user is already a member of the org
	ErrorUserAlreadyInOrg = "the_user_is_already_a_member_of_the_org"

	// org not allow request member
	ErrorOrgNotAllowRequestMember = "org_not_allow_request_member"

	// invalid actor name
	ErrorInvalidActorName = "invalid_actor_name"

	// org fullname is empty
	ErrorOrgFullnameIsEmpty = "org_fullname_is_empty"

	// invalid account
	ErrorInvalidAccount = "invalid_account"

	// invalid org
	ErrorInvalidOrg = "invalid_org"

	// invalid actor
	ErrorInvalidActor = "invalid_actor"

	// invalid user
	ErrorInvalidUser = "invalid_user"

	// invalid requester
	ErrorInvalidRequester = "invalid_requester"
	
	// fullname can't be empty
	ErrorFullnameCanNotBeEmpty = "fullname_can_not_be_empty"

	// failed to update user info
	ErrorFailedToUpdateUserInfo = "failed_to_update_user_info"

	// failed to update git user info
	ErrorFailedToUPdateGitUserInfo = "failed_to_update_git_userinfo"

	// username invalid
	ErrorUsernameInvalid = "username_invalid"
	
	// failed to delete user in git server
	ErrorFailedToDeleteUserInGitServer = "failed_to_delete_user_in_git_server"

	// user already requested to be delete
	ErrorUserAlreadyRequestedToBeDelete = "user_already_requested_to_be_delete"

	// failed to create token
	ErrorFailedToCreateToken = "failed_to_create_token"

	// failed to ecrypt token
	ErrorFailedToEcryptToken = "failed_to_ecrypt_token"

	// input param is empty
	ErrorInputParamIsEmpty = "input_param_is_empty"
	
	// delete token param is empty
	ErrorDeleteTokenParamIsEmpty = "delete_token_param_is_empty"

	// failed to save org
	ErrorFailedToSaveOrg = "failed_to_save_org"

	// nothing changed
	ErrorNothingChanged = "nothing_changed"

	// failed to get org info
	ErrorFailedToGetOrgInfo = "failed_to_get_org_info"

	// failed to get member info
	ErrorFailedToGetMemberInfo = "failed_to_get_member_info"

	// failed to save member for adding member
	ErrorFailedToSaveMemberForAddingMember = "failed_to_save_member_for_adding_member"

	// the only owner can not be removed
	ErrorOnlyOwnerCanNotBeRemoved = "the_only_owner_can_not_be_removed"

	// failed to validate cmd
	ErrorFailedToValidateCmd = "failed_to_validate_cmd"

	// failed to delete git member
	ErrorFailedToDeleteGitMember = "failed_to_delete_git_member"

	// failed to delete member
	ErrorFailedToDeleteMember = "failed_to_delete_member"

	// failed to change owner of org
	ErrorFailedToChangeOwnerOfOrg = "failed_to_change_owner_of_org"
	
	// failed to get members by org name: %s, %s
	ErrorFailedToGetMembersByOrgName = "failed_to_get_members_by_org_name"

	// user %s is already a member of the org %s
	ErrorUserAccountIsAlreadyAMemberOfOrgAccount = "user_account_is_already_a_member_of_the_org_account"

	// failed to add member:%s to org:%s
	ErrorFailedToAddMemberToOrg = "failed_to_add_member"

	// invalid status %s
	ErrorInvalidStatus = "invalid_status"
	
    // user name %s is already taken
	ErrorUsernameIsAlreadyTaken = "user_name_is_already_taken"

	// failed to create platform user: %s
	ErrorFailedToCreatePlatformUser = "failed_to_create_platform_user"

	// failed to save user in db: %s
	ErrorFailToSaveUserInDb = "failed_to_save_user"	
)

// errorImpl
type errorImpl struct {
	code     string
	msg      string
	innerErr error // error info for diagnostic
}

// Error returns the error message.
//
// This function returns the error message of the errorImpl struct.
//
// No parameters.
// Returns a string representing the error message.
func (e errorImpl) Error() string {
	return e.msg
}

// ErrorCode returns the error code.
//
// This function returns the error code of the errorImpl struct.
// The error code is a string representing the type of the error, it could be used for error handling and diagnostic.
//
// No parameters.
// Returns a string representing the error code.
func (e errorImpl) ErrorCode() string {
	return e.code
}

func (e errorImpl) InnerError() error {
	return e.innerErr
}

func InnerErr(err error) error {
	if e, ok := err.(errorImpl); ok {
		return e.innerErr
	}
	return err
}

// New creates a new error with the specified code and message.
func New(code string, msg string, err error) errorImpl {
	v := errorImpl{
		code:     code,
		innerErr: err,
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
func NewNotFound(code string, msg string, err error) notfoudError {
	return notfoudError{errorImpl: New(code, msg, err)}
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
func NewNoPermission(msg string, err error) noPermissionError {
	return noPermissionError{errorImpl: New(errorCodeNoPermission, msg, err)}
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
func NewInvalidParam(msg string, err error) errorImpl {
	return New(errorCodeInvalidParam, msg, err)
}

// NewCountExceeded creates a new error with the specified count exceeded message.
func NewCountExceeded(msg string, err error) errorImpl {
	return New(ErrorCodeCountExceeded, msg, err)
}

// limitRateError
type limitRateError struct {
	errorImpl
}

// OverLimit is a marker method for over limit rate error.
func (l limitRateError) OverLimit() {}

// NewOverLimit creates a new over limit error with the specified code and message.
func NewOverLimit(code string, msg string, err error) limitRateError {
	return limitRateError{errorImpl: New(code, msg, err)}
}

func NewExpired(msg string, err error) errorImpl {
	return New(ErrorCodeExpired, msg, err)
}

// NewCommonRespError creates a new error with the common resp error.
func NewCommonRespError(msg string, err error) errorImpl {
	return New(ErrorBaseCase, msg, err)
}
