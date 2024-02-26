/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides utility functions for handling HTTP errors and error codes.
package controller

import (
	"net/http"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
)

const (
	errorSystemError     = "system_error"
	errorBadRequestBody  = "bad_request_body"
	errorBadRequestParam = "bad_request_param"
)

type errorCode interface {
	ErrorCode() string
}

type errorNotFound interface {
	errorCode

	NotFound()
}

type errorNoPermission interface {
	errorCode

	NoPermission()
}

func httpError(err error) (int, string) {
	if err == nil {
		return http.StatusOK, ""
	}

	sc := http.StatusInternalServerError
	code := errorSystemError

	if v, ok := err.(errorCode); ok {
		code = v.ErrorCode()

		if _, ok := err.(errorNotFound); ok {
			sc = http.StatusNotFound

		} else if _, ok := err.(errorNoPermission); ok {
			sc = http.StatusForbidden

		} else {
			switch code {
			case allerror.ErrorCodeAccessTokenInvalid:
				sc = http.StatusUnauthorized

			case allerror.ErrorCodeLoginIdMissing:
				sc = http.StatusUnauthorized

			case allerror.ErrorCodeLoginIdInvalid:
				sc = http.StatusUnauthorized

			case allerror.ErrorCodeLoginIdNotFound:
				sc = http.StatusUnauthorized

			case allerror.ErrorCodeCSRFTokenMissing:
				sc = http.StatusUnauthorized

			case allerror.ErrorCodeCSRFTokenInvalid:
				sc = http.StatusUnauthorized

			case allerror.ErrorCodeCSRFTokenNotFound:
				sc = http.StatusUnauthorized

			case allerror.ErrorCodeOrgExistModel:
				sc = http.StatusBadRequest

			default:
				sc = http.StatusBadRequest
			}
		}
	}

	return sc, code
}
