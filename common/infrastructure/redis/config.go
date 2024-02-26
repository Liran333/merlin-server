/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package redis provides functionality for interacting with Redis.
package redis

// Config represents the configuration for Redis.
type Config struct {
	IdleSize int    `json:"idle_size"`
	NetWork  string `json:"network"`
	Address  string `json:"address"   required:"true"`
	Password string `json:"password"  required:"true"`
	KeyPair  string `json:"key_pair"`
	DB       int    `json:"db"`
	Timeout  int64  `json:"timeout"`
	DBCert   string `json:"db_cert"`
}

// SetDefault sets the default values for the Config.
func (p *Config) SetDefault() {
	if p.IdleSize <= 0 {
		p.IdleSize = 20
	}

	if p.NetWork == "" {
		p.NetWork = "tcp"
	}

	if p.DB == 0 {
		p.DB = 0
	}

	if p.Timeout == 0 {
		p.Timeout = 10
	}
}
