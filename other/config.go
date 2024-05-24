/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package other provides additional functionality and types.
package other

// Config is a type alias for Config.
type Config struct {
	Analyse Analyse `json:"analyse"`
}

// Analyse is a type alias for Analyse.
type Analyse struct {
	ClientID     string `json:"client_id"     required:"true"`
	ClientSecret string `json:"client_secret" required:"true"`
	GetTokenUrl  string `json:"get_token_url" required:"true"`
}
