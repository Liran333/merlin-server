/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// obs related configuration.
package obs

type Config struct {
	Endpoint  string `json:"endpoint"                  required:"true"`
	AccessKey string `json:"access_key"                required:"true"`
	SecretKey string `json:"secret_key"                required:"true"`
}
