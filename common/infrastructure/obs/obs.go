/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package obs Provide OBS related operations.
package obs

import (
	"io"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

type ObsService interface {
	CreateObject(f io.Reader, bucket, path string) error
	GenFileDownloadURL(bucket, p string, downloadExpiry int) (string, error)
	CopyObject(bucket, dst, src string) error
	DeleteObject(bucket string, path string) error
}

var cli *obs.ObsClient

func Init(cfg *Config) (err error) {
	cli, err = obs.New(cfg.AccessKey, cfg.SecretKey, cfg.Endpoint)
	if err != nil {
		return
	}

	return
}

func Client() *obs.ObsClient {
	return cli
}
