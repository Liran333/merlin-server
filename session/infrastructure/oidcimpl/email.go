/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package oidcimpl provides implementation for OpenID Connect (OIDC) related functionality
// such as sending and verifying email messages.
package oidcimpl

import (
	"bytes"
	"errors"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	libutils "github.com/opensourceways/server-common-lib/utils"
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
		return
	}

	req, err := http.NewRequest(http.MethodPost, impl.getManagerTokenURL, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	var res managerToken
	if err = sendHttpRequest(req, &res); err != nil {
		return
	}

	if res.Status != http.StatusOK {
		err = errors.New("get token error")
		return
	}

	token = res.Token

	return
}

// SendBindEmail sends a bind email request with the provided email and captcha.
func (impl *user) SendBindEmail(email, capt string) (err error) {
	token, err := impl.getManagerToken()
	if err != nil {
		return
	}

	return impl.sendEmail(token, bindEmail, email, capt)
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
	logrus.Errorf("email error: %s", err.Error())

	errinfo := err.Error()
	if strings.Contains(errinfo, infoCodeError) {
		return allerror.New(allerror.ErrorEmailCodeError, "email verify code error", err)
	}

	if strings.Contains(errinfo, infoEmailDuplicateBind1) ||
		strings.Contains(errinfo, infoEmailDuplicateBind2) {
		return allerror.New(allerror.ErrorCodeEmailDuplicateBind, "email duplicate bind", err)
	}

	if strings.Contains(errinfo, infoUserDuplicateBind1) ||
		strings.Contains(errinfo, infoUserDuplicateBind2) {
		return allerror.New(allerror.ErrorCodeUserDuplicateBind, "user duplicate bind", err)
	}

	if strings.Contains(errinfo, infoEmailDuplicateSend) {
		return allerror.New(allerror.ErrorCodeEmailDuplicateSend, "verify code duplicate send", err)
	}

	if strings.Contains(errinfo, infoCodeInvalid) {
		return allerror.New(allerror.ErrorEmailCodeInvalid, "email verify code invalid", err)
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
		return
	}

	req, err := http.NewRequest(http.MethodPost, impl.sendEmailURL, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header.Add("token", token)

	var res normalEmailRes
	err = sendHttpRequest(req, &res)

	if res.Status != http.StatusOK {
		err = errorReturn(err)
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
		return
	}

	req, err := http.NewRequest(http.MethodPost, impl.bindEmailURL, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header.Add("token", token)

	var res normalEmailRes
	err = sendHttpRequest(req, &res)

	if res.Code != http.StatusOK {
		err = errorReturn(err)
	}

	return
}

// VerifyBindEmail verifies the bind email with the provided email, passCode, and userid.
func (impl *user) VerifyBindEmail(email, passCode, userid string) (err error) {
	token, err := impl.getManagerToken()
	if err != nil {
		return
	}

	return impl.verifyBindEmail(token, email, passCode, userid)
}
