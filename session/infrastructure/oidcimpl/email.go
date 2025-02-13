/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package oidcimpl provides implementation for OpenID Connect (OIDC) related functionality
// such as sending and verifying email messages.
package oidcimpl

import (
	"bytes"
	"net/http"
	"strings"

	libutils "github.com/opensourceways/server-common-lib/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
)

const (
	resetPassword = "CHANNEL_RESET_PASSWORD"
	verifyEmail   = "CHANNEL_VERIFY_EMAIL_LINK"
	changeEmail   = "CHANNEL_UPDATE_EMAIL"
	bindEmail     = "CHANNEL_BIND_EMAIL"
	unbindEmail   = "CHANNEL_UNBIND_EMAIL"

	accountTypeEmail = "email"

	infoCodeError           = "E0002"
	infoCodeInvalid         = "E00026"
	infoEmailDuplicateBind1 = "E0004"
	infoEmailDuplicateBind2 = "unique constraint \\\\"
	infoEmailDuplicateSend  = "E00049"
	infoUserDuplicateBind1  = "已绑定邮箱\\\\"
	infoUserDuplicateBind2  = "E00016"
)

type managerBody struct {
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
	GrantType string `json:"grant_type"`
}

type managerToken struct {
	Status int    `json:"status"`
	Token  string `json:"token"`
	Msg    string `json:"msg"`
}

func (impl *user) getManagerToken() (token string, err error) {
	b := managerBody{
		AppId:     impl.cfg.APPId,
		AppSecret: impl.cfg.Secret,
		GrantType: "token",
	}

	body, err := libutils.JsonMarshal(&b)
	if err != nil {
		err = xerrors.Errorf("get manager token error: %w", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, impl.getManagerTokenURL, bytes.NewBuffer(body))
	if err != nil {
		err = xerrors.Errorf("get manager token error: %w", err)
		return
	}

	var res managerToken
	if err = sendHttpRequest(req, &res); err != nil {
		err = xerrors.Errorf("get manager token error: %w", err)
		return
	}

	if res.Status != http.StatusOK {
		err = xerrors.Errorf("get manager token error: status(%d) not ok", res.Status)
		return
	}

	token = res.Token

	return
}

// SendBindEmail sends a bind email request with the provided email and captcha.
func (impl *user) SendBindEmail(email, capt string) (err error) {
	token, err := impl.getManagerToken()
	if err != nil {
		err = xerrors.Errorf("get manager token error: %w", err)
		return
	}

	err = impl.sendEmail(token, bindEmail, email, capt)
	if err != nil {
		return
	}

	return
}

type sendEmail struct {
	Account             string `json:"account"`
	Channel             string `json:"channel"`
	CaptchaVerification string `json:"captchaVerification"`
}

type normalEmailRes struct {
	Code   int `json:"code"`
	Status int `json:"status"`
}

func errorReturn(err error) error {
	errinfo := err.Error()
	if strings.Contains(errinfo, infoCodeInvalid) {
		return allerror.New(allerror.ErrorEmailCodeInvalid, "email verify code invalid", err)
	}

	if strings.Contains(errinfo, infoCodeError) {
		return allerror.New(allerror.ErrorEmailCodeError, "email verify code error", err)
	}

	if strings.Contains(errinfo, infoEmailDuplicateSend) {
		return allerror.New(allerror.ErrorCodeEmailDuplicateSend, "verify code duplicate send", err)
	}

	if strings.Contains(errinfo, infoUserDuplicateBind1) ||
		strings.Contains(errinfo, infoUserDuplicateBind2) {
		return allerror.New(allerror.ErrorCodeUserDuplicateBind, "user duplicate bind", err)
	}

	if strings.Contains(errinfo, infoEmailDuplicateBind1) ||
		strings.Contains(errinfo, infoEmailDuplicateBind2) {
		return allerror.New(allerror.ErrorCodeEmailDuplicateBind, "email duplicate bind", err)
	}

	return allerror.New(allerror.ErrorEmailError, "email bind error", err)
}

func (impl *user) sendEmail(token, channel, email, capt string) (err error) {
	send := sendEmail{
		Account:             email,
		Channel:             channel,
		CaptchaVerification: capt,
	}

	body, err := libutils.JsonMarshal(&send)
	if err != nil {
		err = xerrors.Errorf("marshal error: %w", err)

		return
	}

	req, err := http.NewRequest(http.MethodPost, impl.sendEmailURL, bytes.NewBuffer(body))
	if err != nil {
		err = xerrors.Errorf("new request error: %w", err)

		return
	}

	req.Header.Add("token", token)

	var res normalEmailRes
	err = sendHttpRequest(req, &res)

	if res.Status != http.StatusOK {
		logrus.Errorf("authing email bind err: %v", err)

		err = errorReturn(err)
		logrus.Errorf("error code after parsing authing email bind err: %v", err)

		return
	}

	return
}

type veriEmail struct {
	Account     string `json:"account"`
	Code        string `json:"code"`
	UserId      string `json:"user_id"`
	AccountType string `json:"account_type"`
}

func (impl *user) verifyBindEmail(token, email, passCode, userid string) (err error) {
	veri := veriEmail{
		Account:     email,
		Code:        passCode,
		UserId:      userid,
		AccountType: accountTypeEmail,
	}

	body, err := libutils.JsonMarshal(&veri)
	if err != nil {
		err = xerrors.Errorf("marshal error: %w", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, impl.bindEmailURL, bytes.NewBuffer(body))
	if err != nil {
		err = xerrors.Errorf("new request error: %w", err)
		return
	}

	req.Header.Add("token", token)

	var res normalEmailRes
	err = sendHttpRequest(req, &res)

	if res.Code != http.StatusOK {
		logrus.Errorf("authing email bind err: %v", err)

		err = errorReturn(err)
		logrus.Errorf("error code after parsing authing email bind err: %v", err)

		return
	}

	return
}

// VerifyBindEmail verifies the bind email with the provided email, passCode, and userid.
func (impl *user) VerifyBindEmail(email, passCode, userid string) (err error) {
	token, err := impl.getManagerToken()
	if err != nil {
		err = xerrors.Errorf("get manager token error: %w", err)
		return allerror.New(allerror.ErrorCodeEmailVerifyFailed, err.Error(), err)
	}

	err = impl.verifyBindEmail(token, email, passCode, userid)

	return
}
