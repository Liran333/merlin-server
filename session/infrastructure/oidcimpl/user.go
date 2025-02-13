/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package oidcimpl

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	libutils "github.com/opensourceways/server-common-lib/utils"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/session/domain/repository"
	"github.com/openmerlin/merlin-server/utils"
)

const (
	MaxRetries = 3 // MaxRetries represents the maximum number of retries
	timeout    = 10
)

var userInstance *user

// Config represents the configuration for the authentication service.
type Config struct {
	APPId    string `json:"app_id"        required:"true"`
	Secret   string `json:"secret"        required:"true"`
	Endpoint string `json:"endpoint"      required:"true"`
}

// Init initializes the user instance with the provided configuration.
func Init(v *Config) {
	userInstance = &user{
		cfg:                *v,
		tokenURL:           v.Endpoint + "/oidc/token",
		userInfoURL:        v.Endpoint + "/oidc/user",
		getManagerTokenURL: v.Endpoint + "/manager/token",
		sendEmailURL:       v.Endpoint + "/manager/sendcode",
		bindEmailURL:       v.Endpoint + "/manager/bind/account",
		privacyRevokeUrl:   v.Endpoint + "/manager/privacy/revoke",
	}
}

// NewAuthingUser creates a new instance of the user.
func NewAuthingUser() *user {
	return userInstance
}

type user struct {
	cfg                Config
	tokenURL           string
	userInfoURL        string
	getManagerTokenURL string
	sendEmailURL       string
	bindEmailURL       string
	privacyRevokeUrl   string
}

// GetByAccessToken retrieves user information by access token.
func (impl *user) GetByAccessToken(accessToken string) (userInfo repository.UserInfo, err error) {
	if accessToken == "" {
		err = errors.New("no access token")

		return
	}

	var v struct {
		Name     string `json:"username,omitempty"`
		Picture  string `json:"picture,omitempty"`
		Email    string `json:"email,omitempty"`
		Sub      string `json:"sub,omitempty"`
		FullName string `json:"nickname,omitempty"`
		Phone    string `json:"phone_number,omitempty"`
	}

	if err = impl.getUserInfoByAccessToken(accessToken, &v); err != nil {
		return
	}

	if userInfo.Name, err = primitive.NewAccount(v.Name); err != nil {
		return
	}

	if userInfo.Email, err = primitive.NewUserEmail(v.Email); err != nil {
		return
	}

	if userInfo.AvatarId, err = primitive.NewAvatar(v.Picture); err != nil {
		return
	}

	if userInfo.Fullname, err = primitive.NewAccountFullname(v.FullName); err != nil {
		return
	}

	if userInfo.Phone, err = primitive.NewPhone(v.Phone); err != nil {
		return
	}

	if v.Sub == "" {
		err = errors.New("no sub")

		return
	}
	userInfo.UserId = v.Sub

	return
}

// GetByCode retrieves login information by code and redirectURI.
func (impl *user) GetByCode(code, redirectURI string) (login repository.Login, err error) {
	var v struct {
		AccessToken string `json:"access_token"`
		IdToken     string `json:"id_token"`
	}

	if err = impl.getAccessTokenByCode(code, redirectURI, &v); err != nil {
		return
	}
	defer utils.ClearStringMemory(v.AccessToken)

	if v.IdToken == "" {
		err = errors.New("no id token")

		return
	}

	info, err := impl.GetByAccessToken(v.AccessToken)
	if err == nil {
		login.IDToken = v.IdToken
		login.UserInfo = info
		login.AccessToken = v.AccessToken
	}

	return
}

func (impl *user) getAccessTokenByCode(code, redirectURI string, result interface{}) error {
	body := map[string]string{
		"client_id":     impl.cfg.APPId,
		"client_secret": impl.cfg.Secret,
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  redirectURI,
	}

	value := make(url.Values)
	for k, v := range body {
		value.Add(k, v)
	}

	req, err := http.NewRequest(
		http.MethodPost, impl.tokenURL,
		strings.NewReader(strings.TrimSpace(value.Encode())),
	)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return sendHttpRequest(req, result)
}

func (impl *user) getUserInfoByAccessToken(accessToken string, result interface{}) error {
	req, err := http.NewRequest(http.MethodGet, impl.userInfoURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", accessToken)

	return sendHttpRequest(req, result)
}

// PrivacyRevoke sends a request to revoke privacy for the given user ID.
func (impl *user) PrivacyRevoke(userid string) error {
	var v = struct {
		UserId string `json:"userId"`
	}{
		UserId: userid,
	}

	body, err := libutils.JsonMarshal(&v)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, impl.privacyRevokeUrl, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	token, err := impl.getManagerToken()
	if err != nil {
		return err
	}

	req.Header.Add("token", token)

	return sendHttpRequest(req, nil)
}

func sendHttpRequest(req *http.Request, result interface{}) error {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "merlin-server-authing")
	req.Header.Add("content-type", "application/json")

	httpClient := libutils.NewHttpClient(MaxRetries)
	httpClient.Client.Timeout = time.Duration(timeout) * time.Second

	_, err := httpClient.ForwardTo(req, result)

	return err
}
